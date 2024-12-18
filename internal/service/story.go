package service

import (
	"fmt"
)

// StoryRequest 用于接收客户端发送的 JSON 数据
type StoryRequest struct {
	StoryContent    string `json:"story_content"`
	CharacterChoice string `json:"character_choice"`
	StoryType       string `json:"story_type"`
	ImageType       string `json:"image_type"`
	ChildAgeGroup   string `json:"child_age_group"`
}

// StoryService 是处理故事请求的服务层
type StoryService struct{}

// NewStoryService 创建一个新的 StoryService 实例
func NewStoryService() *StoryService {
	return &StoryService{}
}

// GenerateStory 处理故事请求，并根据请求生成简单的故事内容
func (s *StoryService) GenerateStory(req *StoryRequest) string {
	// 基于请求数据生成一个简单的故事
	return fmt.Sprintf("Here's your story:\n\nStory: %s\nCharacter: %s\nStory Type: %s\nImage Type: %s\nAge Group: %s\n",
		req.StoryContent, req.CharacterChoice, req.StoryType, req.ImageType, req.ChildAgeGroup)
}

// ProcessStoryRequest 用于处理故事请求的业务逻辑
func (s *StoryService) ProcessStoryRequest(req *StoryRequest) string {
	// 在这里可以扩展更多的业务逻辑，比如生成故事、音频等
	story := s.GenerateStory(req)
	// 这里可以加入保存故事到数据库或其他操作

	// 返回处理后的故事内容
	return story
}
