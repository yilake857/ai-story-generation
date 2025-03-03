package service

import (
	"flutterdreams/config"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func LoadConfigForTest(t *testing.T) {
	// 获取当前文件的绝对路径并推算出项目根目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("无法获取当前文件路径")
	}

	// 获取当前文件所在的目录
	projectRoot := filepath.Dir(filepath.Dir(filename))

	// 构造配置文件的绝对路径
	configPath := filepath.Join(projectRoot, "../config", "config.yaml")
	configPath, err := filepath.Abs(configPath) // 转换为绝对路径
	if err != nil {
		t.Fatalf("获取绝对路径失败: %v", err)
	}
	// 加载配置
	_, err = config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}
}

func TestGenerateStory_ValidRequest(t *testing.T) {
	// 调用加载配置的函数
	LoadConfigForTest(t)

	// 创建请求
	req := &StoryRequest{
		StoryContent:    "爱探险的朵拉",
		StoryType:       "冒险",
		ChildAgeGroup:   "3-5岁",
		ImageType:       "卡通风格",
		CharacterChoice: "youxiaoxun",
	}

	resp := &StoryResponse{}

	// 创建 StoryService 实例
	service := StoryService{}

	// 调用 GenerateStory 方法
	err := service.ProcessStoryRequest(req, resp)
	if err != nil {
		t.Errorf("错误：%v", err)
	}
	fmt.Println(resp.StoryContent + resp.ImagePrompt + resp.AudioUrl)
}
