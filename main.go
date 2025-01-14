package main

import (
	"flutterdreams/config"
	"flutterdreams/internal/route"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

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

func main() {
	Hello()
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	router := route.InitRouter()

	// 配置 CORS 允许来自 localhost:3000 的请求
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // 允许的源
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// 使用 CORS 中间件
	handler := c.Handler(router)
	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Server starting at http://%s\n", address)
	err = http.ListenAndServe(address, handler)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
