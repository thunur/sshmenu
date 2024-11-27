package menu

import (
	"fmt"
	"github.com/thunur/sshmenu/internal/config"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func AddServerPrompt() (*config.SSHConfig, error) {
	prompts := []struct {
		label    string
		validate func(string) error
	}{
		{
			label: "服务器名称",
			validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("名称不能为空")
				}
				return nil
			},
		},
		{
			label: "主机地址",
			validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("主机地址不能为空")
				}
				return nil
			},
		},
		{
			label: "用户名",
			validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("用户名不能为空")
				}
				return nil
			},
		},
		{
			label: "端口号",
			validate: func(input string) error {
				port, err := strconv.Atoi(input)
				if err != nil {
					return fmt.Errorf("端口必须是数字")
				}
				if port < 1 || port > 65535 {
					return fmt.Errorf("端口范围必须在 1-65535 之间")
				}
				return nil
			},
		},
	}

	results := make([]string, len(prompts))
	for i, p := range prompts {
		prompt := promptui.Prompt{
			Label:    p.label,
			Validate: p.validate,
		}

		result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("输入取消: %v", err)
		}
		results[i] = result
	}

	// 选择认证方式
	authPrompt := promptui.Select{
		Label: "选择认证方式",
		Items: []string{"SSH密钥", "密码"},
	}

	authIndex, _, err := authPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("选择取消: %v", err)
	}

	port, _ := strconv.Atoi(results[3])
	server := &config.SSHConfig{
		Name: results[0],
		Host: results[1],
		User: results[2],
		Port: port,
	}

	// 根据认证方式获取额外信息
	if authIndex == 0 {
		// SSH密钥
		keyPrompt := promptui.Prompt{
			Label:   "SSH密钥路径",
			Default: "~/.ssh/id_rsa",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("密钥路径不能为空")
				}
				return nil
			},
		}
		keyPath, err := keyPrompt.Run()
		if err != nil {
			return nil, fmt.Errorf("输入取消: %v", err)
		}
		server.KeyPath = strings.Replace(keyPath, "~", os.Getenv("HOME"), 1)
	} else {
		// 密码
		passPrompt := promptui.Prompt{
			Label: "密码",
			Mask:  '*',
		}
		password, err := passPrompt.Run()
		if err != nil {
			return nil, fmt.Errorf("输入取消: %v", err)
		}
		server.Password = password
	}

	return server, nil
}
