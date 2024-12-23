package model

import (
	"flutterdreams/config"
	utils2 "flutterdreams/pkg/utils"
	"flutterdreams/pkg/utils/authv3"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

// 2. 根据故事结果 + character_choice 返回音频文件
func GenerateAudioFromText(q string, voiceName string) (string, error) {
	// 创建请求参数
	paramsMap := map[string][]string{
		"q":         {q},         // 传入的文本
		"voiceName": {voiceName}, // 语音名称
		"format":    {"mp3"},     // 输出格式
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
	fileName, filePath := getTempFilePath()

	// 如果音频生成成功
	if result == nil {
		log.Fatalf("failed to generate audio")
		return "", fmt.Errorf("failed to generate audio")
	}
	// 保存音频文件
	utils2.SaveFile(filePath, result, false)
	log.Printf("save file path: " + filePath)

	return fileName, nil
}

// 获取保存临时文件的路径
func getTempFilePath() (string, string) {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return "", ""
	}
	tempDir := filepath.Join(workingDir, "audio")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return "", ""
	}
	// 生成唯一的文件名
	fileName := uuid.New().String() + ".mp3"
	// 拼接文件的完整路径
	filePath := filepath.Join(tempDir, fileName)
	return fileName, filePath
}
