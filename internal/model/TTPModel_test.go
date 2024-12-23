package model

import (
	"testing"
)

func TestGenerateImage(t *testing.T) {
	userPrompt := "A futuristic city skyline with flying cars"
	GenerateImage(userPrompt)
}
