package route

import (
	"encoding/json"
	"flutterdreams/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
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

	return router
}

// 处理 POST 请求的函数
func CreateStory(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 创建一个 StoryRequest 实例来存储请求体中的数据
	var storyReq service.StoryRequest

	// 解析 JSON 请求体
	err := json.NewDecoder(r.Body).Decode(&storyReq)
	if err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 创建一个新的 StoryService 实例
	storyService := service.NewStoryService()

	// 使用 StoryService 处理故事请求并生成故事
	result := storyService.ProcessStoryRequest(&storyReq)

	// 返回 JSON 响应，表示数据已成功接收
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusCreated)

	// 返回生成的故事内容
	response := map[string]string{
		"status":  "success",
		"message": "Story request received successfully",
		"story":   result,
	}

	// 将响应转换为 JSON 格式
	json.NewEncoder(wr).Encode(response)
}

// 处理 GET 请求的心跳检查
func HealthCheck(wr http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 返回一个简单的 JSON 响应，表示服务正常
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	json.NewEncoder(wr).Encode(map[string]string{"status": "healthy"})
}
