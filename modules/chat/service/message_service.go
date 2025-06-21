package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
)

func (s *ChatService) BuildMessageHistory(ctx context.Context, sessionID int32) ([]model.ChatMessage, error) {
	conversations, err := s.Queries.GetConversationsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var messages []model.ChatMessage
	for _, c := range conversations {
		messages = append(messages, model.ChatMessage{
			Role:    c.Role,
			Content: c.Message,
		})
	}
	return messages, nil
}

func (s *ChatService) SaveUserMessage(ctx context.Context, sessionID int32, content string) error {
	_, err := s.Queries.CreateConversation(ctx, db.CreateConversationParams{
		SessionID: sessionID,
		Message:   content,
		Role:      "user",
	})
	return err
}

func (s *ChatService) SaveAssistantMessage(ctx context.Context, sessionID int32, content string) error {
	_, err := s.Queries.CreateConversation(ctx, db.CreateConversationParams{
		SessionID: sessionID,
		Message:   content,
		Role:      "assistant",
	})
	return err
}

func (s *ChatService) IsProductRelated(ctx context.Context, message string) bool {
	classificationPrompt := []model.ChatMessage{
		{Role: "system", Content: "Apakah ini pertanyaan tentang produk seperti stok, nama produk, harga? Jawab hanya 'ya' atau 'tidak'."},
		{Role: "user", Content: message},
	}

	response, _, err := s.OpenAIService.ChatCompletion(ctx, classificationPrompt)
	if err != nil {
		log.Println("Classification failed:", err)
		return false
	}

	normalized := strings.ToLower(strings.TrimSpace(response.Message))
	return strings.Contains(normalized, "ya")
}

func (s *ChatService) CheckStockByKeyword(ctx context.Context, userMsg string) (string, error) {
	products, err := s.Queries.GetAllProducts(ctx)
	if err != nil {
		return "", err
	}

	for _, p := range products {
		if strings.Contains(strings.ToLower(userMsg), strings.ToLower(p.Name)) {
			return fmt.Sprintf("Stok produk %s tersedia sebanyak %d dengan harga Rp%d", p.Name, p.Stock, p.Price), nil
		}
	}
	return "Maaf, saya tidak menemukan produk yang Anda maksud.", nil
}
