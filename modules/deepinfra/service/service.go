package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"shofy/modules/chat/model"
)

type OpenAIService struct {
	APIKey string
	Model  string
}

func NewOpenAIService(ctx context.Context) *OpenAIService {
	return &OpenAIService{
		APIKey: os.Getenv("DI_API_KEY"),
		Model:  "meta-llama/Llama-4-Maverick-17B-128E-Instruct-FP8", // ganti model sesuai kebutuhan
	}
}

func convertToDeepInfraFormat(messages []model.ChatMessage) []map[string]string {
	var result []map[string]string
	for _, m := range messages {
		result = append(result, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}
	return result
}

func (s *OpenAIService) ChatCompletion(ctx context.Context, messages []model.ChatMessage) (model.ChatResponse, int, error) {
	payload := map[string]interface{}{
		"model":    s.Model,
		"messages": convertToDeepInfraFormat(messages),
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepinfra.com/v1/openai/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return model.ChatResponse{}, http.StatusInternalServerError, err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return model.ChatResponse{}, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return model.ChatResponse{}, resp.StatusCode, fmt.Errorf("DeepInfra error: %s", b)
	}

	var parsed model.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return model.ChatResponse{}, http.StatusInternalServerError, err
	}

	response := model.ChatResponse{
		Message:      parsed.Choices[0].Message.Content,
		FullResponse: parsed,
	}
	return response, http.StatusOK, nil
}
