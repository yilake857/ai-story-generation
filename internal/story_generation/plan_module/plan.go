package plan_module

import (
	"flutterdreams/internal/story_generation/common"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	MAX_ATTEMPTS          = 3
	MAX_SETTING_LENGTH    = 64
	MAX_CHARACTERS        = 3
	MAX_CHARACTERS_LENGTH = 64
	MAX_OUTLINE_SECTIONS  = 5
	MAX_OUTLINE_LENGTH    = 128
)

// PlanInfo 存储故事计划的所有信息
type PlanInfo struct {
	Premise               string
	Setting               string
	Characters            []string
	CharacterStrings      []string
	Outline               string
	OutlineSections       []string
	InferAttributesString string
}

func GeneratePlanInfo(premise string) (*PlanInfo, error) {
	// 初始化 PlanInfo 结构体
	planInfo := &PlanInfo{}

	log.Println("premise: ", premise)
	planInfo.Premise = premise

	// 生成 setting
	setting, err := generateSetting(premise)
	if err != nil {
		return nil, fmt.Errorf("无法生成setting: %v", err)
	}
	log.Println("setting: ", setting)
	planInfo.Setting = setting

	// 生成角色信息
	characters, characterDetails, err := generateCharactersInfos(premise, setting)
	if err != nil {
		return nil, fmt.Errorf("无法生成角色信息: %v", err)
	}
	log.Println("characters: ", characters)
	for _, characterDetail := range characterDetails {
		log.Println("characterDetail: ", characterDetail)
	}
	planInfo.Characters = characters
	planInfo.CharacterStrings = characterDetails

	// 生成 InferAttributesString
	planInfo.InferAttributesString = fmt.Sprintf("前提：%s\n\n背景：%s\n\n角色：\n%s\n\n角色信息：\n%s",
		premise,
		setting,
		strings.Join(characters, "\n"),
		strings.Join(characterDetails, "\n"),
	)

	// 生成故事大纲
	outline, outlineSections, err := generateOutline(planInfo.InferAttributesString)
	if err != nil {
		return nil, err
	}
	log.Println("outline: ", outline)
	log.Println("outlineSections: ", outlineSections)
	for i, section := range outlineSections {
		log.Printf("section%d: %s", i, section)
	}
	planInfo.Outline = outline
	planInfo.OutlineSections = outlineSections

	return planInfo, nil
}

// 生成角色信息
func generateCharactersInfos(premise string, setting string) ([]string, []string, error) {
	// 拼接premise和setting作为前置提醒
	basePrompt := "故事前提: " + premise + "\n\n" + "故事背景: " + setting + "\n\n"

	charactersPrompt := basePrompt + "请生成" + strconv.Itoa(MAX_CHARACTERS) + "个主要角色，" +
		"要求：\n" +
		"1. 用简体中文\n" +
		"2. 不要使用特殊字符、星号或markdown格式\n" +
		"3. 避免使用括号、方括号或任何可能影响文本转语音的符号\n" +
		"4. 每个角色按照1. 2. 3.的格式列出，如 1. 角色名：特点、背景、对故事的影响。\n" +
		"5. 每个角色需要有中文名字（不包含标点符号）和独特的特点背景。\n"
	var characterNames []string
	var characterDetails []string

	for attempts := 0; attempts < MAX_ATTEMPTS; attempts++ {
		charactersBasic, err := common.ChatWithModel(charactersPrompt)
		if err != nil {
			return nil, nil, fmt.Errorf("无法生成角色基本信息: %v", err)
		}

		// 检查是否包含有效的角色信息格式
		if regexp.MustCompile(`^\d+\.\s`).MatchString(charactersBasic) {
			charactersBasic = removeAsterisks(charactersBasic)
			// 解析出角色数量和基本特征
			characterNames = parseCharacterNames(charactersBasic)
			characterDetails = parseCharacterDetails(charactersBasic)

			// 检查是否成功生成角色详细信息
			if len(characterDetails) > 0 {
				break // 如果成功生成，退出循环
			}
		}
	}

	if len(characterDetails) == 0 {
		return nil, nil, fmt.Errorf("未能生成有效的角色信息")
	}
	if len(characterDetails) > MAX_CHARACTERS {
		characterDetails = characterDetails[:MAX_CHARACTERS]
	}
	return characterNames, characterDetails, nil
}

// 解析角色名称
func parseCharacterNames(charactersBasic string) []string {
	var names []string
	lines := strings.Split(charactersBasic, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 查找类似 "1. 角色名" 或 "1. 角色名：" 的模式
		if match, _ := regexp.MatchString(`^\d+\.`, line); match {
			parts := strings.SplitN(line, ".", 2)
			if len(parts) > 1 {
				namePart := strings.TrimSpace(parts[1])
				// 如果有冒号，取冒号前的部分作为名字
				if colonParts := strings.SplitN(namePart, "：", 2); len(colonParts) > 1 {
					names = append(names, strings.TrimSpace(colonParts[0]))
				} else if colonParts := strings.SplitN(namePart, ":", 2); len(colonParts) > 1 {
					names = append(names, strings.TrimSpace(colonParts[0]))
				} else {
					// 否则取第一个词作为名字
					words := strings.Fields(namePart)
					if len(words) > 0 {
						names = append(names, words[0])
					}
				}
			}
		}
	}

	// 清理名字，移除标点符号和强调符号
	for i, name := range names {
		names[i] = cleanChineseName(name)
	}

	return names
}

// 清理中文名字，移除标点符号和强调符号
func cleanChineseName(name string) string {
	// 移除强调符号 **
	name = strings.ReplaceAll(name, "**", "")
	name = strings.ReplaceAll(name, "*", "")

	// 移除常见标点符号
	punctuations := []string{"，", "。", "、", "：", "；", "！", "？", "（", "）", "「", "」", "『", "』", "\"", "\"", "'", "'", "《", "》"}
	for _, p := range punctuations {
		name = strings.ReplaceAll(name, p, "")
	}

	// 移除英文标点符号
	englishPunct := []string{",", ".", ":", ";", "!", "?", "(", ")", "[", "]", "{", "}", "\"", "'", "<", ">"}
	for _, p := range englishPunct {
		name = strings.ReplaceAll(name, p, "")
	}

	return strings.TrimSpace(name)
}

// 解析角色详细信息
func parseCharacterDetails(charactersBasic string) []string {
	var details []string
	lines := strings.Split(charactersBasic, "\n")

	currentDetail := ""
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if match, _ := regexp.MatchString(`^\d+\.`, line); match {
			// 保存之前的详细信息
			if currentDetail != "" {
				details = append(details, currentDetail)
			}
			currentDetail = line
		} else if currentDetail != "" && i < len(lines)-1 {
			// 继续添加到当前详细信息，除非是空行且下一行是新角色
			if line != "" || (i+1 < len(lines) && !regexp.MustCompile(`^\d+\.`).MatchString(lines[i+1])) {
				currentDetail += "\n" + line
			}
		}
	}

	// 添加最后一个详细信息
	if currentDetail != "" {
		details = append(details, currentDetail)
	}

	// 清理详细信息中的强调符号
	for i, detail := range details {
		details[i] = strings.ReplaceAll(detail, "**", "")
		details[i] = strings.ReplaceAll(detail, "*", "")
	}

	return details
}

// 生成故事大纲
func generateOutline(inferAttributesString string) (string, []string, error) {
	var outlineSections []string
	var outlineSectionsRaw string
	var err error

	for i := 0; i < MAX_ATTEMPTS; i++ {
		// 生成故事大纲，明确限制不生成不当内容
		outlinePrompt := fmt.Sprintf("%s\n\n请生成一个完整的第三人称的故事大纲，分为"+strconv.Itoa(MAX_OUTLINE_SECTIONS)+"个主要部分，"+
			"要求："+
			"1. 用简体中文"+
			"2. 不要使用特殊字符、星号或markdown格式"+
			"3. 避免使用括号、方括号或任何可能影响文本转语音的符号"+
			"4. 请确保内容适合所有年龄段，不包含任何不当或敏感的主题"+
			"5. 每个部分之间需要有伏笔响应且有逻辑关系并言简意赅"+
			"输出示例："+
			"1. 大纲1 "+
			"2. 大纲2 ",
			inferAttributesString)

		outlineSectionsRaw, err = common.ChatWithModel(outlinePrompt)
		if err != nil {
			return "", nil, fmt.Errorf("无法生成大纲分段: %v", err)
		}
		// 移除大纲分段中的 * 符号
		outlineSectionsRaw = removeAsterisks(outlineSectionsRaw)
		// 解析大纲分段为数组
		outlineSections = parseOutlineSections(outlineSectionsRaw)

		// 如果 outlineSections 有内容，则退出循环
		if len(outlineSections) > 0 {
			break
		}
	}
	return outlineSectionsRaw, outlineSections, nil // 返回所有部分
}

// 移除字符串中的 * 符号
func removeAsterisks(input string) string {
	return strings.ReplaceAll(input, "*", "")
}

// 解析大纲分段为数组
func parseOutlineSections(outlineSectionsRaw string) []string {
	var sections []string
	lines := strings.Split(outlineSectionsRaw, "\n")

	// 用于存储当前处理的大纲标题
	var currentOutline string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 匹配以下格式：
		// 1. 内容
		// 第1部分:内容
		// 第一部分：内容
		// （第1部分：内容）
		if match, _ := regexp.MatchString(`^(\d+\.|第[一二三四五六七八九十\d]+部分[：:])`, line); match {
			// 如果已经有当前大纲，先保存它
			if currentOutline != "" {
				sections = append(sections, currentOutline)
			}

			// 移除序号和括号，只保留内容
			content := regexp.MustCompile(`^(\d+\.|第[一二三四五六七八九十\d]+部分[：:]|\(|\))`).ReplaceAllString(line, "")
			currentOutline = strings.TrimSpace(content)
		} else if currentOutline != "" && line != "" {
			// 如果这一行不是大纲标题，且不是空行，那么它可能是详细内容的开始
			// 我们只保留大纲标题，不保留详细内容
			// 将当前大纲添加到结果中，然后重置当前大纲
			sections = append(sections, currentOutline)
			currentOutline = ""
		}
	}

	// 处理最后一个大纲
	if currentOutline != "" {
		sections = append(sections, currentOutline)
	}

	return sections
}

// 新增的 generateSetting 函数
func generateSetting(premise string) (string, error) {
	settingPrompt := "故事的前提是: " + premise + "\n\n描述一下故事的背景\n\n" +
		"要求：\n" +
		"1. 用简体中文\n" +
		"2. 不要使用特殊字符、星号或markdown格式\n" +
		"3. 背景要有趣且富有想象力\n" +
		"4. 使用简单明了的语言\n" +
		"5. 避免使用括号、方括号或任何可能影响文本转语音的符号\n" +
		"对故事发展有指导意义" + "\n\n这个故事发生在"
	var setting string
	var err error

	for attempts := 0; attempts < MAX_ATTEMPTS; attempts++ {
		setting, err = common.ChatWithModel(settingPrompt)
		if err == nil && setting != "" {
			//切割setting，只保留前100个字符
			if len(setting) > MAX_SETTING_LENGTH {
				setting = setting[:MAX_SETTING_LENGTH]
			}
			return removeAsterisks(setting), nil
		}
	}
	return "", fmt.Errorf("未能生成有效的 setting")
}
