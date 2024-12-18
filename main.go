package main

import (
	"flutterdreams/config"
	"flutterdreams/internal/route"
	"fmt"
	"log"
	"net/http"
)

// Hello 函数将 ASCII 艺术输出到控制台
func Hello() {
	asciiArt := `
_____.__          __    __                   .___                                     
_/ ____\  |  __ ___/  |__/  |_  ___________  __| _/______   ____ _____    _____   ______
\   __\|  | |  |  \   __\   __\/ __ \_  __ \/ __ |\_  __ \_/ __ \\__  \  /     \ /  ___/
 |  |  |  |_|  |  /|  |  |  | \  ___/|  | \/ /_/ | |  | \/\  ___/ / __ \|  Y Y  \\___ \ 
 |__|  |____/____/ |__|  |__|  \___  >__|  \____ | |__|    \___  >____  /__|_|  /____  >
								   \/           \/             \/     \/      \/     \/ 	
`
	log.Println(asciiArt)
}

// main 函数启动服务器并输出相关信息
func main() {
	// 打印 ASCII 艺术和服务启动信息
	Hello()
	// 读取配置
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// 启动路由
	router := route.InitRouter()
	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	// 设置端口并启动 HTTP 服务
	log.Printf("Server starting at http://%s\n", address)
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
