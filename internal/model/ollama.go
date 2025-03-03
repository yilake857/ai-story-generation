package model

import (
	"context"
	"log"

	"flutterdreams/config"
	"github.com/ollama/ollama/api"
)

func ChatWithOllama(prompt string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	model := config.GetConfig().Ollama.Model // 获取模型
	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		// set streaming to false
		Stream: new(bool),
	}
	var responseContent string
	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		responseContent = resp.Response
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	return responseContent, nil
}
