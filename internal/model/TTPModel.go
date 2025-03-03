package model

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateImage(userPrompt string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	pythonInterpreter := filepath.Join(workingDir, "cmd", "ttp", "venv", "bin", "python3")
	pythonScript := filepath.Join(workingDir, "cmd", "ttp", "generate_image.py")

	cmd := exec.Command(pythonInterpreter, pythonScript)
	cmd.Stdin = bytes.NewBufferString(userPrompt)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run Python script: %v, stderr: %s", err, stderr.String())
	}

	output := out.String()
	if len(output) > 0 {
		log.Printf("Generated image URL: %s", output)
		return output, nil
	} else {
		return "", fmt.Errorf("failed to generate image, no output received")
	}
}
