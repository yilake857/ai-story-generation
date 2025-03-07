//给定premise，生成story

package story_generation

import (
	"flutterdreams/internal/story_generation/draft_module"
	"flutterdreams/internal/story_generation/plan_module"
	"fmt"
	"log"
)

func GenerateStory(premise string) string {
	//plan
	planInfo, err := plan_module.GeneratePlanInfo(premise)
	if err != nil {
		fmt.Println("生成计划信息时出错: ", err)
		return ""
	}

	//Draft
	draft, err := draft_module.GenerateDraft(planInfo.InferAttributesString, planInfo.OutlineSections)
	if err != nil {
		fmt.Println("生成草稿时出错: ", err)
		return ""
	}
	log.Println("draft: ", draft)
	return draft
	//Rewrite

	//Edit

	//Story
}
