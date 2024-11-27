package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/thunur/sshmenu/internal/config"
)

func Connect(cfg *config.SSHConfig) error {
	// 构建 ssh 命令
	args := []string{}

	// 添加端口参数
	if cfg.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", cfg.Port))
	}

	// 添加密钥参数
	if cfg.KeyPath != "" {
		// 处理 ~ 路径
		keyPath := cfg.KeyPath
		if strings.HasPrefix(keyPath, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("获取用户目录失败: %v", err)
			}
			keyPath = strings.Replace(keyPath, "~", home, 1)
		}
		args = append(args, "-i", keyPath)
	}

	// 构建目标地址
	target := fmt.Sprintf("%s@%s", cfg.User, cfg.Host)
	args = append(args, target)

	// 创建命令
	cmd := exec.Command("ssh", args...)

	// 设置标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}
