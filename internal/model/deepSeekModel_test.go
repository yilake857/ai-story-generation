package model

import (
	"flutterdreams/config"
	"fmt"
	"log"
	"testing"
)

// mockConfig 初始化一个模拟的配置
func mockDeepSeekConfig() {
	config.GlobalConfig = config.Config{
		Deepseek: config.DeepseekConfig{
			Api: "sk-da5e66c4edf44f99bf71b32df971839e",
		},
	}
}

func TestChatByDeepSeek(t *testing.T) {
	mockDeepSeekConfig()
	// Example usage
	systemContent := "You are a helpful assistant"
	userContent := "Hello"
	response, err := ChatByDeepSeek(systemContent, userContent)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Response:", response)
}
