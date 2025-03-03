package common

import (
	"flutterdreams/config"
	"log"
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
	configPath := filepath.Join(projectRoot, "../../config", "config.yaml")
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

func TestChatWithModel(t *testing.T) {
	LoadConfigForTest(t)

	// 调用 ChatWithModel 函数
	response, err := ChatWithModel("随机生成一个短故事的前提,字数不超过128个字")
	if err != nil {
		t.Fatalf("调用 ChatWithModel 失败: %v", err)
	}
	log.Println("response: ", response)
	// 检查返回的响应是否为空
	if response == "" {
		t.Error("返回的响应为空")
	}

	// 可以根据需要添加更多的断言来验证响应内容
}
