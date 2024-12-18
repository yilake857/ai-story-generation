package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

// LoadConfig 读取配置文件并解析
func LoadConfig(filename string) (Config, error) {
	var config Config

	// 读取配置文件
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("Error reading config file: %v", err)
	}

	// 解析 YAML 数据
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("Error parsing YAML data: %v", err)
	}

	return config, nil
}
