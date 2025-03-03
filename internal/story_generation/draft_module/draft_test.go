package draft_module

import (
	"flutterdreams/config"
	"log"
	"path/filepath"
	"runtime"
	"strings"
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

func TestGenerateDraft(t *testing.T) {
	// 测试数据
	inferAttributesString := `前提：一个年轻人发现自己可以在梦中控制现实
背景：这个故事发生在一个现代城市，科技发达但人们的生活压力很大
角色：
1. 林宇：年轻程序员，发现自己拥有梦境控制能力
2. 苏瑶：心理咨询师，帮助林宇理解他的能力
角色信息：
1. 林宇：25岁程序员，性格内向但富有创造力，在梦中能改变现实
2. 苏瑶：28岁心理咨询师，善解人意，对超自然现象有研究`

	outlineSections := []string{
		"1. 林宇发现自己能在梦中控制现实，感到兴奋好奇，开始探索。",
		"2. 林宇结识志同道合的苏瑶，两人组成团队共同探索梦境奥秘。",
		"3. 他们的努力引起神秘组织注意，组织企图利用林宇的能力。",
		"4. 林宇和团队面临挑战，努力保护自己和世界安全。",
	}

	// 调用函数
	storyDraft, err := GenerateDraft(inferAttributesString, outlineSections)
	if err != nil {
		t.Fatalf("生成故事草稿失败: %v", err)
	}

	// 验证结果
	if storyDraft == "" {
		t.Error("生成的故事草稿为空")
	}

	// 检查故事草稿是否包含关键角色名称
	if !strings.Contains(storyDraft, "林宇") || !strings.Contains(storyDraft, "苏瑶") {
		t.Error("故事草稿中缺少主要角色名称")
	}

	// 打印生成的故事草稿
	t.Logf("生成的故事草稿:\n%s", storyDraft)
}

func TestGetBestCandidate(t *testing.T) {
	// 测试数据
	inferAttributesString := `前提：一个年轻人发现自己可以在梦中控制现实
背景：这个故事发生在一个现代城市，科技发达但人们的生活压力很大
角色：
1. 林宇：年轻程序员，发现自己拥有梦境控制能力
2. 苏瑶：心理咨询师，帮助林宇理解他的能力`

	// 测试不同位置的段落
	testCases := []struct {
		name    string
		draft   Draft
		wantErr bool
	}{
		{
			name: "测试第一段",
			draft: Draft{
				Index:                 0,
				InferAttributesString: inferAttributesString,
				CurrentSection:        "1. 林宇发现自己能在梦中控制现实，感到兴奋好奇，开始探索。",
				NextOutlineSection:    "2. 林宇结识志同道合的苏瑶，两人组成团队共同探索梦境奥秘。",
			},
			wantErr: false,
		},
		// {
		// 	name: "测试中间段",
		// 	draft: Draft{
		// 		Index:                 1,
		// 		InferAttributesString: inferAttributesString,
		// 		PreOutlineSection:     "1. 林宇发现自己能在梦中控制现实，感到兴奋好奇，开始探索。",
		// 		CurrentSection:        "2. 林宇结识志同道合的苏瑶，两人组成团队共同探索梦境奥秘。",
		// 		NextOutlineSection:    "3. 他们的努力引起神秘组织注意，组织企图利用林宇的能力。",
		// 		PreContent:            "在一个平凡的夜晚，林宇做了一个不平凡的梦...",
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "测试最后一段",
		// 	draft: Draft{
		// 		Index:                 3,
		// 		InferAttributesString: inferAttributesString,
		// 		PreOutlineSection:     "3. 他们的努力引起神秘组织注意，组织企图利用林宇的能力。",
		// 		CurrentSection:        "4. 林宇和团队面临挑战，努力保护自己和世界安全。",
		// 		PreContent:            "神秘组织的人开始对林宇展开行动...",
		// 	},
		// 	wantErr: false,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 获取最佳候选集
			bestCandidate, err := getBestCandidate(tc.draft)

			// 检查错误
			if (err != nil) != tc.wantErr {
				t.Errorf("getBestCandidate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// 验证结果不为空
			if bestCandidate == "" {
				t.Error("getBestCandidate() 返回空字符串")
			}

			// 检查返回的内容是否包含关键角色名称
			if !strings.Contains(bestCandidate, "林宇") {
				t.Error("生成的内容中缺少主角'林宇'")
			}

			// 根据段落位置检查特定内容
			switch tc.draft.Index {
			case 0:
				if !strings.Contains(strings.ToLower(bestCandidate), "开始") {
					t.Error("第一段应该包含故事的开始")
				}
			case 3:
				if !strings.Contains(strings.ToLower(bestCandidate), "结局") && !strings.Contains(strings.ToLower(bestCandidate), "最终") {
					t.Error("最后一段应该包含故事的结局")
				}
			}

			// 打印生成的内容
			t.Logf("生成的内容:\n%s", bestCandidate)
		})
	}
}
