package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type SSHConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	KeyPath  string `json:"key_path,omitempty" mapstructure:"key_path"`
	Password string `json:"password,omitempty" mapstructure:"password"`
}

type Config struct {
	Servers []SSHConfig `json:"servers"`
}

func InitConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".sshmenu")
	configPath := filepath.Join(configDir, "config.json")

	// 检查配置目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %v", err)
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建默认空配置
		defaultConfig := Config{
			Servers: []SSHConfig{},
		}
		viper.SetConfigType("json")
		viper.Set("servers", defaultConfig.Servers)
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("创建配置文件失败: %v", err)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	return nil
}

func LoadConfig() ([]SSHConfig, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}
	return config.Servers, nil
}

func AddServer(server SSHConfig) error {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 检查是否存在同名配置
	for _, s := range config.Servers {
		if s.Name == server.Name {
			return fmt.Errorf("服务器名称 '%s' 已存在", server.Name)
		}
	}

	config.Servers = append(config.Servers, server)
	viper.Set("servers", config.Servers)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	return nil
}

func DeleteServer(name string) error {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 查找并删除指定服务器
	found := false
	newServers := make([]SSHConfig, 0)
	for _, server := range config.Servers {
		if server.Name == name {
			found = true
			continue
		}
		newServers = append(newServers, server)
	}

	if !found {
		return fmt.Errorf("未找到名为 '%s' 的服务器", name)
	}

	// 更新配置
	config.Servers = newServers
	viper.Set("servers", config.Servers)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	return nil
}
