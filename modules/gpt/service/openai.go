package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"shofy/modules/chat/model"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	Model  string
	Client openai.Client
}

const (
	ModelR1Turbo = "deepseek-ai/DeepSeek-R1-Turbo"
)

func NewOpenAIService(ctx context.Context) *OpenAIService {
	endpoint := os.Getenv("DEEPINFRA_URL")
	if endpoint == "" {
		log.Fatal("DEEPINFRA_URL is not set")
	}

	secretAPIKey := os.Getenv("DEEPINFRA_API_KEY")
	if secretAPIKey == "" {
		log.Fatal("DEEPINFRA_API_KEY is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(secretAPIKey),
		option.WithBaseURL(endpoint),
	)
	return &OpenAIService{
		Model:  ModelR1Turbo,
		Client: client,
	}
}

func (s *OpenAIService) ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) (model.ChatResponse, int, error) {
	chatCompletion, err := s.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    s.Model,
	})
	if err != nil {
		return model.ChatResponse{}, http.StatusInternalServerError, err
	}
	chatResponse := model.ChatResponse{
		Message: chatCompletion.Choices[0].Message.Content,
	}
	return chatResponse, http.StatusOK, nil
}
