package model

import (
	"fmt"
	"testing"
)

func TestChatWithOllama(t *testing.T) {
	prompt := "为什么天空是蓝色的？"

	response, err := ChatWithOllama(prompt)
	if err != nil {
		t.Fatalf("ChatWithOllama() 返回错误: %v", err)
	}
	if response == "" {
		t.Error("ChatWithOllama() 返回的内容为空")
	}
	fmt.Println(response)
}
