package draft_module

import (
	"flutterdreams/internal/story_generation/common"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Draft 结构体用于保存大纲段落的上下文信息
type Draft struct {
	Index                 int
	InferAttributesString string

	PreOutlineSection  string
	CurrentSection     string
	NextOutlineSection string

	PreContent string
	Content    string
}

const (
	MAX_CANDIDATE_SIZE = 1
)

// 返回：故事草稿
func GenerateDraft(InferAttributesString string, outlineSections []string) (string, error) {
	story_draft := ""
	sectionsCount := len(outlineSections)
	Drafts := make([]Draft, sectionsCount)
	// 遍历大纲段落
	for i, currentSection := range outlineSections {
		draft := Draft{
			Index:                 i,
			CurrentSection:        currentSection,
			InferAttributesString: InferAttributesString,
		}

		// 只有非第一段才设置 PreOutlineSection
		if i > 0 {
			draft.PreOutlineSection = outlineSections[i-1]
			draft.PreContent = Drafts[i-1].Content
		}

		// 只有非最后一段才设置 NextOutlineSection
		if i < sectionsCount-1 {
			draft.NextOutlineSection = outlineSections[i+1]
		}

		//取出最理想的候选集
		content, err := getBestCandidate(draft)
		if err != nil {
			return "", fmt.Errorf("无法生成候选集: %v", err)
		}
		story_draft += content
		Drafts[i] = draft
	}
	return story_draft, nil
}

func getBestCandidate(draft Draft) (string, error) {
	prompt := construct_prompt(draft)
	//生成max_candidate_size个候选集
	candidateList := make([]string, MAX_CANDIDATE_SIZE)
	//分数0-10
	currentScore := 0
	bestCandidate := ""
	for i := 0; i < MAX_CANDIDATE_SIZE; i++ {
		candidate, err := generateCandidate(prompt)
		if err != nil {
			return "", fmt.Errorf("无法生成候选集: %v", err)
		}
		log.Println("Draft Index: ", draft.Index, " candidate Index: ", i, " candidate: ", candidate)
		candidateList[i] = candidate
		score := getScore(draft, candidateList[i])
		if score >= currentScore {
			currentScore = score
			bestCandidate = candidateList[i]
		}
	}
	//对bestCandidate去掉多余的符号 写个函数
	bestCandidate = removeExtraSymbols(bestCandidate)
	log.Println("bestCandidate: ", bestCandidate)
	return bestCandidate, nil
}

func removeExtraSymbols(candidate string) string {
	// 去掉被**包裹的内容及**符号
	candidate = regexp.MustCompile(`\*\*[^*]*\*\*`).ReplaceAllString(candidate, "")

	// 去掉序号（如 "1.", "2."）和第一个冒号前的"大纲"字样
	candidate = regexp.MustCompile(`^(\d+\.\s*|.*?大纲[： :])(.+)$`).ReplaceAllString(candidate, "$2")

	// 去掉多余的空白字符
	candidate = strings.TrimSpace(candidate)
	// 将多个空格替换为单个空格
	candidate = regexp.MustCompile(`\s+`).ReplaceAllString(candidate, " ")

	return candidate
}

// 生成单个候选集
func generateCandidate(prompt string) (string, error) {
	candidate, err := common.ChatWithModel(prompt)
	if err != nil {
		return "", fmt.Errorf("无法生成候选集: %v", err)
	}
	return candidate, nil
}

// 获取候选集的分数
func getScore(draft Draft, candidate string) int {
	return 0
}

func construct_prompt(draft Draft) string {
	//draft 是第一段
	if draft.PreOutlineSection == "" {
		prompt := fmt.Sprintf("背景信息：%s\n当前段落大纲：%s\n下一段落大纲：%s\n\n开头（即当前）段落的全文如下\n",
			draft.InferAttributesString,
			draft.CurrentSection,
			draft.NextOutlineSection,
		)
		return prompt
	}
	//draft 是最后一段
	if draft.NextOutlineSection == "" {
		prompt := fmt.Sprintf("背景信息：%s\n前一段大纲：%s\n前一段内容：%s\n当前一段大纲：\n%s\n故事结尾全文如下\n",
			draft.InferAttributesString,
			draft.PreOutlineSection,
			draft.PreContent,
			draft.CurrentSection,
		)
		return prompt
	}
	//draft 是中间段
	prompt := fmt.Sprintf("背景信息：%s\n前一段落大纲：%s\n前一段落内容：%s\n下一段落大纲：%s\n当前段落大纲：\n%s\n当前段落的全文如下\n",
		draft.InferAttributesString,
		draft.PreOutlineSection,
		draft.PreContent,
		draft.NextOutlineSection,
		draft.CurrentSection,
	)
	return prompt
}
