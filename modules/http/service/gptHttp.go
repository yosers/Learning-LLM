package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GPTService struct {
	BaseService *BaseService
	Model       string
}

const (
	ModelR1Turbo = "deepseek-ai/DeepSeek-R1-Turbo"
)

func NewGPTService(ctx context.Context, dbPool *pgxpool.Pool, http *http.Client, queries *db.Queries) *GPTService {
	baseUrl := os.Getenv("DEEPINFRA_URL")
	if baseUrl == "" {
		log.Fatal("openapi url is not set")
	}
	openAPIKey := os.Getenv("DEEPINFRA_API_KEY")
	if openAPIKey == "" {
		log.Fatal("openapi key is not set")
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", openAPIKey),
	}
	baseService := NewBaseService(baseUrl)
	baseService.Headers = headers
	return &GPTService{
		BaseService: baseService,
		Model:       ModelR1Turbo,
	}
}

func (s *GPTService) CreateChat(ctx context.Context, request map[string]interface{}) (*http.Response, int, error) {
	requestParams := map[string]interface{}{
		"model": s.Model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": request["message"],
			},
		},
	}
	requestBody, err := json.Marshal(requestParams)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	response, status, err := s.BaseService.Post(ctx, "chat/completions", requestBody)
	if err != nil {
		return nil, status, err
	}

	return response, http.StatusOK, nil
}
