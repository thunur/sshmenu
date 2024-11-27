package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thunur/sshmenu/internal/config"
	"github.com/thunur/sshmenu/internal/menu"
	"github.com/thunur/sshmenu/internal/ssh"

	"github.com/manifoldco/promptui"
)

func main() {
	// 定义命令行参数
	isManage := flag.Bool("m", false, "进入管理模式")
	flag.Parse()

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 加载现有配置
	configs, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 如果没有配置，直接进入添加流程
	if len(configs) == 0 {
		fmt.Println("未发现任何服务器配置，请先添加服务器。")
		if err := handleAddServer(); err != nil {
			fmt.Fprintf(os.Stderr, "添加服务器失败: %v\n", err)
			os.Exit(1)
		}
		// 重新加载配置
		configs, err = config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "重新加载配置失败: %v\n", err)
			os.Exit(1)
		}
	}

	if *isManage {
		// 管理模式
		if err := manageMode(); err != nil {
			fmt.Fprintf(os.Stderr, "管理操作失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 直接进入连接模式
		if err := handleConnect(configs); err != nil {
			fmt.Fprintf(os.Stderr, "连接失败: %v\n", err)
			os.Exit(1)
		}
	}
}

func manageMode() error {
	for {
		actions := []string{"添加服务器", "删除服务器", "返回"}
		prompt := promptui.Select{
			Label: "选择操作",
			Items: actions,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("选择失败: %v", err)
		}

		switch index {
		case 0: // 添加服务器
			if err := handleAddServer(); err != nil {
				fmt.Printf("添加服务器失败: %v\n", err)
				continue
			}
		case 1: // 删除服务器
			if err := handleDeleteServer(); err != nil {
				fmt.Printf("删除服务器失败: %v\n", err)
				continue
			}
		case 2: // 返回
			return nil
		}
	}
}

func handleDeleteServer() error {
	configs, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if len(configs) == 0 {
		return fmt.Errorf("没有可删除的服务器")
	}

	selected, err := menu.ShowMenu(configs, menu.MenuTypeDelete)
	if err != nil {
		return fmt.Errorf("选择服务器失败: %v", err)
	}

	// 确认删除
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("确认删除服务器 '%s' (yes/no)", selected.Name),
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("操作取消")
	}

	if result != "yes" && result != "y" {
		return fmt.Errorf("操作取消")
	}

	if err := config.DeleteServer(selected.Name); err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}

	fmt.Printf("成功删除服务器 '%s'\n", selected.Name)
	return nil
}

func handleConnect(configs []config.SSHConfig) error {
	selected, err := menu.ShowMenu(configs, menu.MenuTypeConnect)
	if err != nil {
		return fmt.Errorf("选择服务器失败: %v", err)
	}

	return ssh.Connect(selected)
}

func handleAddServer() error {
	server, err := menu.AddServerPrompt()
	if err != nil {
		return err
	}

	if err := config.AddServer(*server); err != nil {
		return err
	}

	fmt.Printf("成功添加服务器 '%s'\n", server.Name)
	return nil
}
