package route

import (
	"encoding/json"
	"flutterdreams/internal/service"
	"flutterdreams/internal/story_generation"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// 初始化路由
func InitRouter() *httprouter.Router {
	router := httprouter.New()

	// 设置 POST 路由
	router.POST("/story", CreateStory)

	// 设置 GET 路由用于心跳检查
	router.GET("/health", HealthCheck)

	// 获取音频文件
	router.GET("/getAudio", GetAudio)
	// 生成故事
	router.POST("/generateStory", GenerateStory)
	return router
}

func CreateStory(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 创建一个 StoryRequest 实例来存储请求体中的数据
	var storyReq service.StoryRequest
	var storyResp service.StoryResponse

	// 解析 JSON 请求体
	err := json.NewDecoder(r.Body).Decode(&storyReq)
	if err != nil {
		// 这里不再调用 http.Error 直接返回，改为返回错误信息
		logError(wr, "Invalid request body", err)
		return
	}

	// 创建一个新的 StoryService 实例
	storyService := service.NewStoryService()

	// 使用 StoryService 处理故事请求并生成故事
	err = storyService.ProcessStoryRequest(&storyReq, &storyResp)
	if err != nil {
		logError(wr, "Failed to process story request", err)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK) // 设置 HTTP 状态码为 200 OK

	response := map[string]interface{}{
		"status":       "success",
		"message":      "Story request received successfully",
		"title":        storyResp.StoryTitle,
		"story":        storyResp.StoryContent,
		"image_prompt": storyResp.ImagePrompt,
		"audio_url":    storyResp.AudioUrl,
		"image_url":    strings.ReplaceAll(storyResp.ImageUrl, "\n", ""),
	}

	// 将响应转换为 JSON 格式并返回
	err = json.NewEncoder(wr).Encode(response)
	if err != nil {
		logError(wr, "Error encoding response", err)
	}
}

// 处理 GET 请求的心跳检查
func HealthCheck(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 返回一个简单的 JSON 响应，表示服务正常
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	json.NewEncoder(wr).Encode(map[string]string{"status": "healthy"})
}

func GetAudio(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 获取请求中的文件名参数
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		logError(wr, "Missing filename parameter", fmt.Errorf("filename is required"))
		return
	}

	// 获取音频文件路径
	audioFilePath := getAudioFilePath(fileName)

	// 检查文件是否存在
	if _, err := os.Stat(audioFilePath); os.IsNotExist(err) {
		logError(wr, "File not found", err)
		return
	}

	// 设置响应头以告知浏览器音频文件类型
	wr.Header().Set("Content-Type", "audio/mpeg")

	// 设置响应头以指定音频文件下载
	wr.Header().Set("Content-Disposition", "inline; filename=\""+fileName+"\"")

	// 直接调用 ServeFile 不需要错误捕获，因为它直接写入响应流
	http.ServeFile(wr, r, audioFilePath)
}

// 获取音频文件的路径，根据文件名来确定
func getAudioFilePath(fileName string) string {
	// 假设音频文件存储在项目根目录下的 "audio" 目录
	// 且文件名必须以 ".mp3" 后缀结尾
	if !strings.HasSuffix(fileName, ".mp3") {
		fileName += ".mp3" // 默认添加 .mp3 后缀
	}
	return filepath.Join("audio", fileName)
}

// 记录日志并返回错误信息
func logError(wr http.ResponseWriter, message string, err error) {
	log.Printf("%s: %v", message, err)
	http.Error(wr, fmt.Sprintf("%s: %v", message, err), http.StatusInternalServerError)
}

// StoryGenerateRequest 定义请求体结构
type StoryGenerateRequest struct {
	Premise string `json:"premise"`
}

// StoryGenerateResponse 定义响应体结构
type StoryGenerateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GenerateStory(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 解析请求体
	var req StoryGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(wr, "Invalid request body", err)
		return
	}

	// 验证premise不为空
	if req.Premise == "" {
		logError(wr, "Missing premise", fmt.Errorf("premise is required"))
		return
	}

	// 调用 plan_module 生成故事计划
	story_generation.GenerateStory(req.Premise)

	// 构造响应
	response := StoryGenerateResponse{
		Status:  "success",
		Message: "Story plan generated successfully",
	}

	// 设置响应头
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)

	// 编码并发送响应
	if err := json.NewEncoder(wr).Encode(response); err != nil {
		logError(wr, "Error encoding response", err)
		return
	}
}
