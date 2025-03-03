package plan_module

import (
	"flutterdreams/config"
	"log"
	"path/filepath"
	"runtime"
	"testing"
)

// 执行测试函数前先加载配置文件
func init() {
	LoadConfigForTest()
}

func LoadConfigForTest() {
	// 获取当前文件的绝对路径并推算出项目根目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("无法获取当前文件路径")
	}

	// 获取当前文件所在的目录
	projectRoot := filepath.Dir(filepath.Dir(filename))

	// 构造配置文件的绝对路径
	configPath := filepath.Join(projectRoot, "../../config", "config.yaml")
	configPath, err := filepath.Abs(configPath) // 转换为绝对路径
	if err != nil {
		log.Fatalf("获取绝对路径失败: %v", err)
	}
	// 加载配置
	_, err = config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
}

func TestGenerateCharactersInfos(t *testing.T) {
	// 测试用的前提和背景
	premise := "一个年轻人发现自己可以在梦中控制现实"
	setting := "这个故事发生在一个现代城市，科技发达但人们的生活压力很大"

	// 调用函数
	characterNames, characterDetails, err := generateCharactersInfos(premise, setting)
	if err != nil {
		t.Fatalf("GenerateCharactersInfos失败: %v", err)
	}

	log.Println("characterNames: ", characterNames)
	log.Println("characterDetails: ", characterDetails)
}

func TestGenerateOutline(t *testing.T) {
	inferAttributesString := "前提：一个年轻的女孩在森林中迷路了。\n\n背景：这个故事发生在一个神秘的森林，充满了奇幻的生物。\n\n角色：\n1. 小红：勇敢的女孩，善于解决问题。\n2. 狼：狡猾的生物，试图引导小红走向危险。"

	outline, outlineSections, err := generateOutline(inferAttributesString)
	if err != nil {
		t.Fatalf("GenerateOutline失败: %v", err)
	}

	log.Println("outline: ", outline)
	log.Println("outlineSections: ", outlineSections)
}

func TestGeneratePlanInfo(t *testing.T) {
	t.Parallel() // 允许并行执行此测试

	premise := "一个年轻人发现自己可以在梦中控制现实"
	planInfo, err := GeneratePlanInfo(premise)
	if err != nil {
		t.Fatalf("生成计划信息时出错: %v", err)
	}

	log.Println("planInfo: ", planInfo)
}
