# SSHMenu - SSH Connection Manager
# SSH 连接管理器

A terminal-based SSH connection manager that helps you manage and connect to your servers easily.
一个基于终端的 SSH 连接管理器，帮助您轻松管理和连接服务器。

## Features 功能特点

- Simple and intuitive terminal UI 简单直观的终端界面
- JSON configuration file support 支持 JSON 配置文件
- Password and SSH key authentication 支持密码和 SSH 密钥认证
- Server management (add/delete) 服务器管理（添加/删除）
- Interactive server selection 交互式服务器选择

## Installation 安装

```bash
Clone the repository 克隆仓库
git clone https://github.com/thunur/sshmenu.git
```
### Build 编译
```
cd sshmenu
go build -o sshmenu main.go
```

```bash
./sshmenu
```

```bash
go install github.com/thunur/sshmenu@latest
```

This will directly show the server selection menu.
这将直接显示服务器选择菜单。

### Management Mode 管理模式

```bash
./sshmenu -m
``` 

Enter management mode to add or delete servers.
进入管理模式以添加或删除服务器。

## Configuration 配置

The configuration file is stored at `~/.sshmenu/config.json`
配置文件存储在 `~/.sshmenu/config.json`

Example configuration 配置示例:


```json
{
  "servers": [
    {
      "name": "Production Server",
      "host": "prod.example.com",
      "user": "admin",
      "port": 22,
      "key_path": "~/.ssh/id_rsa"
    },
    {
      "name": "Test Server",
      "host": "test.example.com",
      "user": "developer",
      "port": 22,
      "password": "your-password"
    }
  ]
}
```

## Key Bindings 快捷键

- `↑/↓`: Navigate through items 上下移动选项
- `Enter`: Select item 选择项目
- `Ctrl+C`: Exit 退出
- `Ctrl+D`: Exit 退出

## Dependencies 依赖

- github.com/spf13/viper
- github.com/manifoldco/promptui
- golang.org/x/crypto/ssh


## FAQ 常见问题

**Q: How to backup my configuration? 如何备份配置？**  
A: Simply copy the `~/.sshmenu/config.json` file.
直接复制 `~/.sshmenu/config.json` 文件即可。

**Q: How to use SSH key authentication? 如何使用 SSH 密钥认证？**  
A: When adding a new server, choose "SSH Key" as authentication method and provide the path to your private key.
添加新服务器时，选择"SSH 密钥"作为认证方式，并提供私钥路径。

**Q: Can I edit the configuration file manually? 可以手动编辑配置文件吗？**  
A: Yes, you can edit `~/.sshmenu/config.json` directly. Just make sure to follow the correct