package model

import (
	"flutterdreams/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

// mockConfig 初始化一个模拟的配置
func mockConfig() {
	config.GlobalConfig = config.Config{
		DoubaoConfig: config.DoubaoConfig{
			Api: "20f5c3b8-cfa8-4cfc-b24b-b0678e7175b2",
		},
	}
}

func TestGenerateStory(t *testing.T) {
	// 初始化模拟配置
	mockConfig()

	// 定义测试用例
	tests := []struct {
		name           string
		systemContent  string
		userContent    string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Valid input",
			systemContent:  "你是一名讲故事专家",
			userContent:    "讲一个关于小红帽的故事",
			expectedOutput: "从前有个小红帽，她住在森林边缘...",
			expectError:    false,
		},
		{
			name:          "Invalid API Key",
			systemContent: "你是一名讲故事专家",
			userContent:   "讲一个关于勇敢兔子的故事",
			expectError:   true,
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := GenerateStory(tt.systemContent, tt.userContent)

			// 检查是否期望错误
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}
