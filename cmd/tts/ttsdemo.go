package main

import (
    utils2 "flutterdreams/pkg/utils"
	"flutterdreams/pkg/utils/authv3"
)

// 您的应用ID
var appKey = "1edbbd24cea9eefe"

// 您的应用密钥
var appSecret = "MVoWYIpTulSid8Ta1MsS00WtRH3wsRRd"

// 合成音频保存路径, 例windows路径：PATH = "C:\\tts\\media.mp3"
var path = "/Users/zhaoyu/go/ai-story-generation/media.mp3"

func main() {
	// 添加请求参数
	paramsMap := createRequestParams()
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	result := utils2.DoPost("https://openapi.youdao.com/ttsapi", header, paramsMap, "audio")
	// 打印返回结果
	if result != nil {
		utils2.SaveFile(path, result, false)
		print("save file path: " + path)
	}

}

func createRequestParams() map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/%E8%AF%AD%E9%9F%B3%E5%90%88%E6%88%90TTS/API%E6%96%87%E6%A1%A3/%E8%AF%AD%E9%9F%B3%E5%90%88%E6%88%90%E6%9C%8D%E5%8A%A1/%E8%AF%AD%E9%9F%B3%E5%90%88%E6%88%90%E6%9C%8D%E5%8A%A1-API%E6%96%87%E6%A1%A3.html
	*/
	q := "金玲瑶是世界上最好看的女人"
	voiceName := "youxiaoxun"
	format := "mp3"

	return map[string][]string{
		"q":         {q},
		"voiceName": {voiceName},
		"format":    {format},
	}
}
