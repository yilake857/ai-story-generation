// 重写模块
package rewrite_module

import (
	"flutterdreams/internal/story_generation/common"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// 使用 common.Draft 替代导入 draft_module.Draft
type Draft = common.Draft

// 打分维度常量
const (
	COHERENCE_WEIGHT     = 0.4  // 连贯性权重
	QUALITY_WEIGHT       = 0.3  // 内容质量权重
	FLUENCY_WEIGHT       = 0.3  // 表达流畅度权重
	MAX_SCORE_PER_ASPECT = 10.0 // 每个维度的最高分
)

// 对候选集打分
func GetScore(draft Draft, candidate string) (float64, error) {
	// 获取连贯性分数
	coherenceScore, err := scoreCoherence(draft, candidate)
	if err != nil {
		return 0, fmt.Errorf("连贯性打分失败: %v", err)
	}

	// 获取内容质量分数
	qualityScore, err := scoreQuality(draft, candidate)
	if err != nil {
		return 0, fmt.Errorf("内容质量打分失败: %v", err)
	}

	// 获取表达流畅度分数
	fluencyScore, err := scoreFluency(draft, candidate)
	if err != nil {
		return 0, fmt.Errorf("表达流畅度打分失败: %v", err)
	}

	// 计算加权总分
	totalScore := calculateTotalScore(coherenceScore, qualityScore, fluencyScore)

	log.Printf("Draft Index: %d, 连贯性: %.1f, 内容质量: %.1f, 表达流畅度: %.1f, 总分: %.1f",
		draft.Index, coherenceScore, qualityScore, fluencyScore, totalScore)

	return totalScore, nil
}

// 评估连贯性
func scoreCoherence(draft Draft, candidate string) (float64, error) {
	prompt := constructCoherencePrompt(draft, candidate)
	response, err := common.ChatWithModel(prompt)
	if err != nil {
		return 0, err
	}

	score, err := extractScore(response, "连贯性")
	if err != nil {
		return 0, err
	}

	return score, nil
}

// 评估内容质量
func scoreQuality(draft Draft, candidate string) (float64, error) {
	prompt := constructQualityPrompt(draft, candidate)
	response, err := common.ChatWithModel(prompt)
	if err != nil {
		return 0, err
	}

	score, err := extractScore(response, "内容质量")
	if err != nil {
		return 0, err
	}

	return score, nil
}

// 评估表达流畅度
func scoreFluency(draft Draft, candidate string) (float64, error) {
	prompt := constructFluencyPrompt(draft, candidate)
	response, err := common.ChatWithModel(prompt)
	if err != nil {
		return 0, err
	}

	score, err := extractScore(response, "表达流畅度")
	if err != nil {
		return 0, err
	}

	return score, nil
}

// 构建连贯性评分提示
func constructCoherencePrompt(draft Draft, candidate string) string {
	var builder strings.Builder

	builder.WriteString("请作为一位专业的文学评论家，评估以下故事段落的连贯性，给出1.0-10.0分的评分（可以包含小数点后一位）。\n\n")

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

	// 添加待评分的段落
	builder.WriteString("待评分段落：\n")
	builder.WriteString(candidate)
	builder.WriteString("\n\n")

	// 添加评分标准
	builder.WriteString("连贯性评分标准（1.0-10.0分）：\n")
	builder.WriteString("- 9.0-10.0分：段落内部逻辑完美，与前后文衔接自然，情节发展顺畅，伏笔和呼应恰到好处\n")
	builder.WriteString("- 7.0-8.9分：段落内部逻辑清晰，与前后文衔接良好，情节发展基本顺畅，有一定的伏笔和呼应\n")
	builder.WriteString("- 5.0-6.9分：段落内部逻辑基本清晰，与前后文有一定衔接，情节发展有些跳跃，伏笔和呼应不够明显\n")
	builder.WriteString("- 3.0-4.9分：段落内部逻辑不够清晰，与前后文衔接不够紧密，情节发展较为跳跃，缺乏伏笔和呼应\n")
	builder.WriteString("- 1.0-2.9分：段落内部逻辑混乱，与前后文衔接生硬，情节发展断裂，没有伏笔和呼应\n\n")

	builder.WriteString("请仔细分析后给出评分，只输出一个数字作为评分结果，格式如下：\n")
	builder.WriteString("连贯性：X.X\n")

	return builder.String()
}

// 构建内容质量评分提示
func constructQualityPrompt(draft Draft, candidate string) string {
	var builder strings.Builder

	builder.WriteString("请作为一位专业的文学评论家，评估以下故事段落的内容质量，给出1.0-10.0分的评分（可以包含小数点后一位）。\n\n")

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

	// 添加待评分的段落
	builder.WriteString("待评分段落：\n")
	builder.WriteString(candidate)
	builder.WriteString("\n\n")

	// 添加评分标准
	builder.WriteString("内容质量评分标准（1.0-10.0分）：\n")
	builder.WriteString("- 9.0-10.0分：内容丰富深刻，主题表达清晰有力，情节设计巧妙，人物形象鲜明，细节描写生动\n")
	builder.WriteString("- 7.0-8.9分：内容较为丰富，主题表达清晰，情节设计合理，人物形象较为鲜明，细节描写较好\n")
	builder.WriteString("- 5.0-6.9分：内容基本充实，主题表达基本清晰，情节设计基本合理，人物形象和细节描写一般\n")
	builder.WriteString("- 3.0-4.9分：内容较为单薄，主题表达不够清晰，情节设计平淡，人物形象模糊，细节描写不足\n")
	builder.WriteString("- 1.0-2.9分：内容空洞，主题表达混乱，情节设计不合理，人物形象扁平，几乎没有细节描写\n\n")

	builder.WriteString("请仔细分析后给出评分，只输出一个数字作为评分结果，格式如下：\n")
	builder.WriteString("内容质量：X.X\n")

	return builder.String()
}

// 构建表达流畅度评分提示
func constructFluencyPrompt(draft Draft, candidate string) string {
	var builder strings.Builder

	builder.WriteString("请作为一位专业的文学评论家，评估以下故事段落的表达流畅度，给出1.0-10.0分的评分（可以包含小数点后一位）。\n\n")

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

	// 添加待评分的段落
	builder.WriteString("待评分段落：\n")
	builder.WriteString(candidate)
	builder.WriteString("\n\n")

	// 添加评分标准
	builder.WriteString("表达流畅度评分标准（1.0-10.0分）：\n")
	builder.WriteString("- 9.0-10.0分：语言优美流畅，句式多样灵活，用词精准丰富，修辞手法恰当，节奏感强\n")
	builder.WriteString("- 7.0-8.9分：语言流畅，句式较为多样，用词准确，有一定修辞手法，节奏感较好\n")
	builder.WriteString("- 5.0-6.9分：语言基本流畅，句式变化一般，用词基本准确，修辞手法和节奏感一般\n")
	builder.WriteString("- 3.0-4.9分：语言不够流畅，句式单一，用词不够准确，几乎没有修辞手法，节奏感差\n")
	builder.WriteString("- 1.0-2.9分：语言生硬，句式混乱，用词不当，没有修辞手法，缺乏节奏感\n\n")

	builder.WriteString("请仔细分析后给出评分，只输出一个数字作为评分结果，格式如下：\n")
	builder.WriteString("表达流畅度：X.X\n")

	return builder.String()
}

// 从响应中提取分数
func extractScore(response string, dimension string) (float64, error) {
	regex := regexp.MustCompile(dimension + `：(\d+\.\d+|\d+)`)
	match := regex.FindStringSubmatch(response)

	if len(match) < 2 {
		return 0, fmt.Errorf("无法从响应中提取%s分数", dimension)
	}

	score, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, fmt.Errorf("%s分数解析失败: %v", dimension, err)
	}

	return clampScore(score), nil
}

// 限制分数在1.0-10.0范围内
func clampScore(score float64) float64 {
	if score < 1.0 {
		return 1.0
	}
	if score > MAX_SCORE_PER_ASPECT {
		return MAX_SCORE_PER_ASPECT
	}
	return score
}

// 计算加权总分
func calculateTotalScore(coherenceScore, qualityScore, fluencyScore float64) float64 {
	weightedScore := coherenceScore*COHERENCE_WEIGHT +
		qualityScore*QUALITY_WEIGHT +
		fluencyScore*FLUENCY_WEIGHT

	return weightedScore
}
