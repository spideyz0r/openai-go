package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	openaiEndpoint = "https://api.openai.com/v1/"
)

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey}
}

type CompletionRequest struct {
	Message     []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices []struct {
		Message      Message `json:"message"`
		Index        int     `json:"index"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

func (c *OpenAIClient) GetCompletion(model string, messages []Message, temperature float32) (CompletionResponse, error) {
	var completionResponse CompletionResponse
	url := fmt.Sprintf("%s%s", openaiEndpoint, "chat/completions")

	requestBody := &CompletionRequest{
		Model:       model,
		Message:     messages,
		Temperature: temperature,
	}

	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return completionResponse, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return completionResponse, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return completionResponse, err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), &completionResponse)
	return completionResponse, nil
}
