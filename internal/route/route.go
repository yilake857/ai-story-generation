package route

import (
	"encoding/json"
	"flutterdreams/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 定义结构体来接收请求的参数
type StoryRequest struct {
	StoryContent    string `json:"story_content"`
	CharacterChoice string `json:"character_choice"`
	StoryType       string `json:"story_type"`
	ImageType       string `json:"image_type"`
	ChildAgeGroup   string `json:"child_age_group"`
}

// 初始化路由
func InitRouter() *httprouter.Router {
	router := httprouter.New()

	// 设置 POST 路由
	router.POST("/story", CreateStory)

	// 设置 GET 路由用于心跳检查
	router.GET("/health", HealthCheck)

	//获取音频文件
	router.GET("/getAudio", GetAudio)
	return router
}

func CreateStory(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 创建一个 StoryRequest 实例来存储请求体中的数据
	var storyReq service.StoryRequest
	var storyResp service.StoryResponse

	// 解析 JSON 请求体
	err := json.NewDecoder(r.Body).Decode(&storyReq)
	if err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 创建一个新的 StoryService 实例
	storyService := service.NewStoryService()

	// 使用 StoryService 处理故事请求并生成故事
	err = storyService.ProcessStoryRequest(&storyReq, &storyResp)
	if err != nil {
		http.Error(wr, "Failed to process story request: "+err.Error(), http.StatusInternalServerError) // 内部服务器错误
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK) // 设置 HTTP 状态码为 200 OK

	response := map[string]interface{}{
		"status":       "success",
		"message":      "Story request received successfully",
		"story":        storyResp.StoryContent, // 返回故事内容
		"image_prompt": storyResp.ImagePrompt,  // 如果有图片提示
		"audio_base64": storyResp.AudioBase64,  // 如果有音频内容
	}

	// 将响应转换为 JSON 格式并返回
	err = json.NewEncoder(wr).Encode(response)
	if err != nil {
		http.Error(wr, "Error encoding response: "+err.Error(), http.StatusInternalServerError) // 返回编码错误
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
		http.Error(wr, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	// 获取音频文件路径
	audioFilePath := getAudioFilePath(fileName)

	// 检查文件是否存在
	if _, err := os.Stat(audioFilePath); os.IsNotExist(err) {
		http.Error(wr, "File not found", http.StatusNotFound)
		return
	}

	// 设置响应头以告知浏览器音频文件类型
	wr.Header().Set("Content-Type", "audio/mpeg")

	// 设置响应头以指定音频文件下载
	wr.Header().Set("Content-Disposition", "inline; filename=\""+fileName+"\"")

	// 读取并返回文件内容
	http.ServeFile(wr, r, audioFilePath)
}

// 获取音频文件的路径，根据文件名来确定
func getAudioFilePath(fileName string) string {
	// 假设音频文件存储在项目根目录下的 "audio" 目录
	// 且文件名必须以 ".mp3" 后缀结尾
	if !strings.HasSuffix(fileName, ".mp3") {
		fileName += ".mp3" // 默认添加 .mp3 后缀
	}
	return filepath.Join("temp", fileName)
}
