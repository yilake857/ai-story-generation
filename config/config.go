package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DoubaoConfig struct {
	Api string `yaml:"api"`
}

type YoudaoTTSConfig struct {
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
	Model     string `yaml:"model"`
}

type DeepseekConfig struct {
	Api   string `yaml:"api"`
	Model string `yaml:"model"`
}

type OllamaConfig struct {
	Model string `yaml:"model"`
}

type Config struct {
	Server       ServerConfig    `yaml:"server"`
	DoubaoConfig DoubaoConfig    `yaml:"doubao"`
	YoudaoTTS    YoudaoTTSConfig `yaml:"youdaoTTS"`
	Deepseek     DeepseekConfig  `yaml:"deepseek"`
	Ollama       OllamaConfig    `yaml:"ollama"`
	DefaultModel string          `yaml:"default_model"`
}

var (
	GlobalConfig Config
	once         sync.Once // 确保配置只被加载一次
)

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

	// 使用 sync.Once 确保全局配置只加载一次
	once.Do(func() {
		GlobalConfig = config
	})

	return config, nil
}

// GetConfig 获取全局配置
func GetConfig() Config {
	return GlobalConfig
}
