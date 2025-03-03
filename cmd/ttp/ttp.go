package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// 输出当前工作目录
	fmt.Println("Current directory:", dir)

	userPrompt := "A futuristic city skyline with flying cars"
	cmd := exec.Command("./cmd/ttp/venv/bin/python3", "./cmd/ttp/generate_image.py") // "generate_image.py" 是你的 Python 脚本路径
	cmd.Stdin = bytes.NewBufferString(userPrompt)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run Python script: %v, stderr: %s", err, stderr.String())
	}
	output := out.String()

	if len(output) > 0 {
		fmt.Println("Generated image URL:", output)
	} else {
		fmt.Println("Failed to generate image.")
	}
}
