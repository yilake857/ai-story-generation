package model

import (
	"bytes"
	"encoding/json"
	"flutterdreams/config"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
	Error   struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// chatByDeepSeek sends a chat request to DeepSeek and returns the response message.
func ChatByDeepSeek(systemContent string, userContent string) (string, error) {
	apiKey := config.GetConfig().Deepseek.Api
	url := "https://api.deepseek.com/chat/completions"

	// Prepare the request body
	chatRequest := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Stream: false,
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set the headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Print the full API response (for debugging purposes)
	fmt.Println("API Response:", string(body))

	// Check for API errors in the response
	if resp.StatusCode != 200 {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return "", fmt.Errorf("failed to parse error response: %v", err)
		}
		return "", fmt.Errorf("API error: %s (code: %s)", errorResponse.Error.Message, errorResponse.Error.Code)
	}

	// Parse the JSON response
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Check if we have any choices in the response
	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices found in the response")
}
