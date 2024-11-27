package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thunur/sshmenu/internal/config"
	"golang.org/x/crypto/ssh"
)

func Connect(cfg *config.SSHConfig) error {
	var authMethods []ssh.AuthMethod

	// 处理密码认证
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	// 处理密钥认证
	if cfg.KeyPath != "" {
		// 展开 ~ 到用户主目录
		if strings.HasPrefix(cfg.KeyPath, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("获取用户主目录失败: %v", err)
			}
			cfg.KeyPath = filepath.Join(homeDir, cfg.KeyPath[1:])
		}

		key, err := os.ReadFile(cfg.KeyPath)
		if err != nil {
			return fmt.Errorf("读取密钥文件失败 (%s): %v", cfg.KeyPath, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("解析密钥失败: %v", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("未配置任何认证方式 (需要配置密码或密钥)")
	}

	config := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：在生产环境中应该验证主机密钥
		BannerCallback: func(message string) error {
			fmt.Println(message)
			return nil
		},
	}

	// 连接服务器
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), config)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %v", err)
	}
	defer session.Close()

	// 设置终端模式
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// 请求伪终端
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		return fmt.Errorf("请求伪终端失败: %v", err)
	}

	// 设置标准输入输出
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 启动shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("启动shell失败: %v", err)
	}

	// 等待会话结束
	if err := session.Wait(); err != nil {
		return fmt.Errorf("会话异常结束: %v", err)
	}

	return nil
}
