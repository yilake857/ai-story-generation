package main

import (
	"context"
	"fmt"
	ark "github.com/sashabaranov/go-openai"
)

func main() {
	config := ark.DefaultConfig("20f5c3b8-cfa8-4cfc-b24b-b0678e7175b2")
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	client := ark.NewClientWithConfig(config)

	fmt.Println("----- standard request -----")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model: "ep-20241213180423-txjtj",
			Messages: []ark.ChatCompletionMessage{
				{
					Role:    ark.ChatMessageRoleSystem,
					Content: "你是一名讲故事专家",
				},
				{
					Role:    ark.ChatMessageRoleUser,
					Content: "我想听有关小红帽的故事",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
