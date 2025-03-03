package common

import (
	"flutterdreams/config"
	"flutterdreams/internal/model"
	"fmt"
)

// ChatWithModel 根据配置文件选择模型并调用相应的生成函数
func ChatWithModel(userContent string) (string, error) {
	defaultModel := config.GetConfig().DefaultModel // 从配置文件读取默认模型

	switch defaultModel {
	case "doubao":
		return model.GenerateStory("", userContent) // 调用 doubao 模型
	case "ollama":
		return model.ChatWithOllama(userContent) // 调用 ollama 模型
	case "deepseek":
		return model.ChatByDeepSeek("", userContent) // 调用 deepseek 模型
	default:
		return "", fmt.Errorf("未定义的模型: %s", defaultModel)
	}
}
