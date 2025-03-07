package edit_module

import (
	"flutterdreams/internal/story_generation/common"
	"fmt"
	"strings"
)

// Rewrite 函数用于检查并修正故事中的事实一致性错误
// 参数：
// - draft: 当前段落的上下文信息
// - candidate: 需要检查和修正的文本内容
// 返回：
// - 修正后的文本内容
func Rewrite(draft common.Draft, candidate string) (string, error) {
	// 构建提示词
	prompt := constructRewritePrompt(draft, candidate)

	// 调用模型进行修正
	response, err := common.ChatWithModel(prompt)
	if err != nil {
		return "", fmt.Errorf("调用模型修正文本失败: %v", err)
	}

	// 清理响应文本，去除可能的前缀说明
	cleanedResponse := cleanResponse(response)

	return cleanedResponse, nil
}

// 构建用于修正文本的提示词
func constructRewritePrompt(draft common.Draft, candidate string) string {
	var builder strings.Builder

	builder.WriteString("请作为一位专业的文学编辑，检查并修正以下故事段落中可能存在的事实一致性错误。\n\n")

	// 添加背景信息
	builder.WriteString("背景信息：\n")
	builder.WriteString(draft.InferAttributesString)
	builder.WriteString("\n\n")

	// 添加上下文信息
	if draft.PreOutlineSection != "" {
		builder.WriteString("前一段大纲：")
		builder.WriteString(draft.PreOutlineSection)
		builder.WriteString("\n")

		builder.WriteString("前一段内容：")
		builder.WriteString(draft.PreContent)
		builder.WriteString("\n\n")
	}

	builder.WriteString("当前段落大纲：")
	builder.WriteString(draft.CurrentSection)
	builder.WriteString("\n\n")

	if draft.NextOutlineSection != "" {
		builder.WriteString("下一段大纲：")
		builder.WriteString(draft.NextOutlineSection)
		builder.WriteString("\n\n")
	}

	// 添加需要修正的段落
	builder.WriteString("需要修正的段落：\n")
	builder.WriteString(candidate)
	builder.WriteString("\n\n")

	// 添加修正要求
	builder.WriteString("修正要求：\n")
	builder.WriteString("1. 检查并修正段落中与背景信息或前文内容不一致的地方\n")
	builder.WriteString("2. 确保人物名称、地点、事件等细节前后一致\n")
	builder.WriteString("3. 修正逻辑矛盾或时间线错误\n")
	builder.WriteString("4. 保持原文的风格和语气\n")
	builder.WriteString("5. 不要添加新的情节，只修正事实一致性问题\n")
	builder.WriteString("6. 如果没有发现问题，请直接返回原文\n\n")

	builder.WriteString("请直接返回修正后的完整段落，不要包含解释或说明。\n")

	return builder.String()
}

// 清理模型返回的响应，去除可能的前缀说明
func cleanResponse(response string) string {
	// 去除可能的"修正后的段落："等前缀
	prefixes := []string{
		"修正后的段落：",
		"修正后的文本：",
		"修改后的段落：",
		"修改后：",
		"修正后：",
	}

	cleanedResponse := response
	for _, prefix := range prefixes {
		if strings.HasPrefix(cleanedResponse, prefix) {
			cleanedResponse = cleanedResponse[len(prefix):]
			break
		}
	}

	return strings.TrimSpace(cleanedResponse)
}
