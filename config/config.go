package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/nzai/go-utility/io"
	"github.com/nzai/go-utility/path"
)

const (
	configFile = "project.json"
)

type Config struct {
	ServerAddress string
}

//	当前系统配置
var configValue *Config = nil

//	设置配置文件
func ReadConfig() error {

	root, err := path.GetStartupDir()
	if err != nil {
		return fmt.Errorf("获取起始目录失败:%s", err.Error())
	}

	//	构造配置文件路径
	filePath := filepath.Join(root, configFile)
	if !io.IsExists(filePath) {
		return fmt.Errorf("配置文件%s不存在", filePath)
	}

	//	读取文件
	buffer, err := io.ReadAllBytes(filePath)
	if err != nil {
		return err
	}

	//	解析配置项
	configValue = &Config{}
	err = json.Unmarshal(buffer, configValue)
	if err != nil {
		return err
	}

	if configValue == nil {
		return fmt.Errorf("配置文件错误")
	}

	return nil
}

//	获取当前系统配置
func Get() *Config {
	return configValue
}
