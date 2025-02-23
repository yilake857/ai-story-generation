package service

import (
	"flutterdreams/config"
	"flutterdreams/internal/model"
	"fmt"
	"log"
	"strings"
)

// 用于接收客户端发送的 JSON 数据
type StoryRequest struct {
	StoryContent    string `json:"story_content"`
	CharacterChoice string `json:"character_choice"`
	StoryType       string `json:"story_type"`
	ImageType       string `json:"image_type"`
	ChildAgeGroup   string `json:"child_age_group"`
}

type StoryResponse struct {
	StoryTitle   string `json:"story_title"`
	StoryContent string `json:"story_content"`
	ImagePrompt  string `json:"image_prompt"`
	AudioUrl     string `json:"audio_url"`
	ImageUrl     string `json:"image_url"`
}

// 是处理故事请求的服务层
type StoryService struct{}

// 创建一个新的 StoryService 实例
func NewStoryService() *StoryService {
	return &StoryService{}
}

// 用于处理故事请求的业务逻辑 逻辑线路
func (s *StoryService) ProcessStoryRequest(req *StoryRequest, resp *StoryResponse) error {
	// 1. story_content + story_type + child_age_group 生成提示词，返回故事结果
	err := s.GenerateStory(req, resp)
	if err != nil {
		log.Printf("生成故事时发生错误: %v", err) // 记录错误，但不停止执行
	}

	// 2. 根据故事结果 + character_choice 返回音频文件
	err = generateAudioFromText(req, resp)
	if err != nil {
		log.Printf("生成音频时发生错误: %v", err) // 记录错误，但不停止执行
	}

	// 3. 根据图片提示词 + image_type 返回图片文件
	err = generateImageFromText(req, resp)
	if err != nil {
		log.Printf("生成图片时发生错误: %v", err) // 记录错误，但不停止执行
	}

	return nil // 不返回错误，确保处理继续进行
}

// 1. story_content + story_type + child_age_group 生成提示词，返回故事结果
func (s *StoryService) GenerateStory(req *StoryRequest, resp *StoryResponse) error {
	// 检查输入是否有效
	if req.StoryContent == "" || req.StoryType == "" || req.ChildAgeGroup == "" || req.ImageType == "" {
		return fmt.Errorf("请求参数无效，请提供故事内容、故事类型和儿童年龄组")
	}

	// 生成故事内容的提示词
	storyPrompt := fmt.Sprintf(
		"根据以下内容生成一个有趣的儿童故事,去掉不相关的内容只输出故事,回复文字的UTF-8编码长度不能超过2000且文字和符号总共不能超过600!"+
			"\n故事的主题：%s\n故事的类型：%s\n儿童的年龄段：%s\n"+
			"你需要按照如下格式进行回复\n"+
			"故事题目：...\n"+
			"故事内容:...",
		req.StoryContent,
		req.StoryType,
		req.ChildAgeGroup,
	)
	log.Printf("storyPrompt:%s", storyPrompt)
	// 调用模型生成故事内容
	storyContent, err := model.GenerateStory(
		"你是一名故事生成的专家，请根据以下提示生成一个有趣的故事。",
		storyPrompt,
	)
	if err != nil {
		log.Printf("生成故事内容时发生错误: %v", err)
		return fmt.Errorf("生成故事内容时发生错误: %v", err)
	}
	//log.Printf("storyContent:%s", storyContent)
	//处理故事题目和故事内容
	title, story := extractStoryInfo(storyContent)
	resp.StoryTitle = title
	resp.StoryContent = story
	log.Printf("StoryTitle:%s", title)
	log.Printf("StoryContent:%s", story)

	// 生成图片提示词的提示词
	imagePromptInput := fmt.Sprintf(
		"根据这个故事内容主旨生成供生成图片的提示词，图片类型：%s\n故事内容：%s",
		req.ImageType,
		storyContent,
	)
	// 调用模型生成图片提示词
	imagePrompt, err := model.GenerateStory(
		"你是一名生成故事的专家，请根据以下提示生成一个适合图片生成的提示词。",
		imagePromptInput,
	)
	if err != nil {
		log.Printf("生成图片提示词时发生错误: %v", err)
		return fmt.Errorf("生成图片提示词时发生错误: %v", err)
	}
	log.Printf("imagePrompt:%s", imagePrompt)
	resp.ImagePrompt = imagePrompt
	return nil
}

// 2. 根据故事结果 + character_choice 返回音频文件
func generateAudioFromText(req *StoryRequest, resp *StoryResponse) error {
	// 检查输入是否有效
	if resp.StoryContent == "" || req.CharacterChoice == "" {
		return fmt.Errorf("请求参数无效，请提供故事内容、音频角色信息")
	}

	if len(resp.StoryContent) > 2048 {
		resp.StoryContent = resp.StoryContent[:2048] // 截断到前 2048 个字符
		log.Printf("StoryContent exceeded 2048 characters, truncated to: %s", resp.StoryContent)
	}

	// 调用生成音频的函数
	fileName, err := model.GenerateAudioFromText(resp.StoryContent, req.CharacterChoice)
	if err != nil {
		log.Printf("Failed to generate audio: %v", err)
		return err
	}

	// 构建音频文件的访问 URL
	address := fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)
	resp.AudioUrl = "http://" + address + "/getAudio?filename=" + fileName
	log.Printf("Audio URL: %s", resp.AudioUrl)

	return nil
}

// 3. 根据图片提示词 + image_type 返回图片文件
func generateImageFromText(req *StoryRequest, resp *StoryResponse) error {
	if resp.ImagePrompt == "" {
		err := fmt.Errorf("ImagePrompt is empty")
		log.Printf("Error: %v", err)
		return err
	}
	imageUrl, err := model.GenerateImage(resp.ImagePrompt)
	if err != nil {
		log.Printf("Failed to generate image: %v", err)
		return err
	}
	resp.ImageUrl = imageUrl
	return nil
}

// 提取故事的题目和内容
func extractStoryInfo(story string) (string, string) {
	// 找到故事中的第一行（题目）和剩余内容（故事内容）
	lines := strings.SplitN(story, "\n", 2) // 按照第一个换行符分割

	// 如果格式正确，返回题目和内容
	if len(lines) == 2 {
		// 去掉 "故事题目：" 和 "故事内容：" 前缀
		title := strings.Replace(lines[0], "故事题目：", "", 1)
		content := strings.Replace(lines[1], "故事内容：", "", 1)
		return title, content
	}

	// 如果格式不对，返回默认题目和内容
	return "生成的故事", story
}
