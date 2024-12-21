package service

import (
	"encoding/base64"
	"flutterdreams/config"
	"flutterdreams/internal/model"
	utils2 "flutterdreams/pkg/utils"
	"flutterdreams/pkg/utils/authv3"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
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
	StoryContent string `json:"story_content"`
	ImagePrompt  string `json:"image_prompt"`
	AudioBase64  string `json:"audio_base64"` //当前直接传到前端 之后再保存下来或者。。。
}

// 是处理故事请求的服务层
type StoryService struct{}

// 创建一个新的 StoryService 实例
func NewStoryService() *StoryService {
	return &StoryService{}
}

// 用于处理故事请求的业务逻辑 逻辑线路
func (s *StoryService) ProcessStoryRequest(req *StoryRequest, resp *StoryResponse) error {
	//1. story_content + story_type + child_age_group 生成提示词，返回故事结果
	err := s.GenerateStory(req, resp)
	if err != nil {
		// 返回空字符串和错误
		return fmt.Errorf("failed to generate story: %w", err)
	}

	//2. 根据故事结果 + character_choice 返回音频文件
	err = generateAudioFromText(req, resp)
	if err != nil {
		// 返回空字符串和错误
		return fmt.Errorf("failed to generate audio: %w", err)
	}

	//3. 根据图片提示词 + image_type 返回图片文件
	return err
}

// 1. story_content + story_type + child_age_group 生成提示词，返回故事结果
func (s *StoryService) GenerateStory(req *StoryRequest, resp *StoryResponse) error {
	// 检查输入是否有效
	if req.StoryContent == "" || req.StoryType == "" || req.ChildAgeGroup == "" || req.ImageType == "" {
		return fmt.Errorf("请求参数无效，请提供故事内容、故事类型和儿童年龄组")
	}

	// 生成故事内容的提示词
	storyPrompt := fmt.Sprintf(
		"根据以下内容生成一个有趣的儿童故事,回复文字的UTF-8编码长度不能超过2048!\n故事的主题：%s\n故事的类型：%s\n儿童的年龄段：%s\n",
		req.StoryContent,
		req.StoryType,
		req.ChildAgeGroup,
	)
	fmt.Printf("%s", storyPrompt)
	// 调用模型生成故事内容
	storyContent, err := model.GenerateStory(
		"你是一名儿童故事专家，请根据以下提示生成一个有趣的故事。",
		storyPrompt,
	)
	if err != nil {
		return fmt.Errorf("生成故事内容时发生错误: %v", err)
	}
	resp.StoryContent = storyContent

	// 生成图片提示词的提示词
	imagePromptInput := fmt.Sprintf(
		"根据这个故事内容主旨生成供生成图片的提示词，图片类型：%s\n故事内容：%s",
		req.ImageType,
		storyContent,
	)

	// 调用模型生成图片提示词
	imagePrompt, err := model.GenerateStory(
		"你是一名儿童故事专家，请根据以下提示生成一个适合图片生成的提示词。",
		imagePromptInput,
	)
	if err != nil {
		return fmt.Errorf("生成图片提示词时发生错误: %v", err)
	}
	resp.ImagePrompt = imagePrompt

	// 返回生成的故事和图片提示词
	return nil
}

// 2. 根据故事结果 + character_choice 返回音频文件
func generateAudioFromText(req *StoryRequest, resp *StoryResponse) error {
	// 检查输入是否有效
	if req.StoryContent == "" || req.CharacterChoice == "" {
		return fmt.Errorf("请求参数无效，请提供故事内容、音频角色信息")
	}

	// 创建请求参数
	paramsMap := map[string][]string{
		"q":         {resp.StoryContent},   // 传入的文本
		"voiceName": {req.CharacterChoice}, // 语音名称
		"format":    {"mp3"},               // 输出格式
	}

	// 设置请求头
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	// 添加鉴权参数
	authv3.AddAuthParams(config.GetConfig().YoudaoTTS.AppKey, config.GetConfig().YoudaoTTS.AppSecret, paramsMap)

	// 请求api服务并获取生成的音频文件数据
	result := utils2.DoPost("https://openapi.youdao.com/ttsapi", header, paramsMap, "audio")

	// 获取保存的文件路径
	path := getTempFilePath()

	// 如果音频生成成功
	if result != nil {
		// 保存音频文件
		utils2.SaveFile(path, result, false)
		print("save file path: " + path)

		//// 将音频文件编码为 Base64
		//encodedAudio, err := encodeFileToBase64(path)
		//if err != nil {
		//	// 错误时返回错误信息
		//	return fmt.Errorf("Error encoding file to base64: %v", err)
		//}
		//
		//// 将 Base64 编码后的音频数据赋值给响应对象
		//resp.AudioBase64 = encodedAudio
		return nil
	}

	// 如果没有返回结果，返回错误
	return fmt.Errorf("failed to generate audio")
}

// 获取保存临时文件的路径
func getTempFilePath() string {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}
	workingDir = filepath.Dir(filepath.Dir(workingDir))
	tempDir := filepath.Join(workingDir, "temp")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return ""
	}
	fileName := uuid.New().String() + ".mp3"
	filePath := filepath.Join(tempDir, fileName)
	return filePath
}

// 将音频文件编码为 Base64 字符串
func encodeFileToBase64(filePath string) (string, error) {
	// 读取音频文件
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 将文件数据转换为 Base64 编码
	encoded := base64.StdEncoding.EncodeToString(fileData)
	return encoded, nil
}
