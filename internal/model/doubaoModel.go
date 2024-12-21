package model

import (
	"context"
	"flutterdreams/config"
	"fmt"
	ark "github.com/sashabaranov/go-openai"
)

func GenerateStory(systemContent string, userContent string) (string, error) {
	//读取配置文件
	config := ark.DefaultConfig(config.GetConfig().DoubaoConfig.Api)
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	client := ark.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model: "ep-20241213180423-txjtj",
			Messages: []ark.ChatCompletionMessage{
				{
					Role:    ark.ChatMessageRoleSystem,
					Content: systemContent,
				},
				{
					Role:    ark.ChatMessageRoleUser,
					Content: userContent,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}
	return resp.Choices[0].Message.Content, nil
}
