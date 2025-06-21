package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
	"strings"

	deepinfraService "shofy/modules/deepinfra/service"

	utils "shofy/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatService struct {
	DBPool  *pgxpool.Pool
	Queries *db.Queries

	// GPTService *service.GPTService
	// AzureOpenAI *azureService.AzureOpenAI
	OpenAIService *deepinfraService.OpenAIService
}

func NewChatService(ctx context.Context, dbPool *pgxpool.Pool, queries *db.Queries) *ChatService {
	openaiService := deepinfraService.NewOpenAIService(ctx)
	return &ChatService{
		DBPool:        dbPool,
		Queries:       queries,
		OpenAIService: openaiService,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, chat model.ChatPayload, channelID int) (model.ChatResponse, int, error) {

	session, err := s.Queries.GetCurrentSessions(ctx, db.GetCurrentSessionsParams{
		ChannelID: int32(channelID),
		UserID:    2,
	})

	var hasNoSession bool = false
	messages := []model.ChatMessage{}
	if err != nil {
		log.Println("WHATSHAP err 2: " + err.Error())
		session, err = s.Queries.CreateSession(ctx, db.CreateSessionParams{
			ChannelID: int32(channelID),
			UserID:    2,
		})
		if err != nil {
			log.Println("WHATSHAP err 3: " + err.Error())
			return model.ChatResponse{}, http.StatusInternalServerError, err
		}
		hasNoSession = true
	} else {
		conversations, err := s.Queries.GetConversationsBySessionID(ctx, session.ID)
		if err != nil {
			hasNoSession = true
		}
		for _, conversation := range conversations {
			messages = append(messages, model.ChatMessage{Role: conversation.Role, Content: conversation.Message})
		}
	}

	conversation, err := s.Queries.CreateConversation(ctx, db.CreateConversationParams{
		SessionID: session.ID,
		Message:   chat.Message,
		Role:      "user",
	})

	if hasNoSession {
		messages = append(messages, model.ChatMessage{Role: "user", Content: conversation.Message})
	}

	messages = append(messages, model.ChatMessage{Role: "user", Content: conversation.Message})

	chatResponse, status, err := s.OpenAIService.ChatCompletion(ctx, messages)
	if err != nil {
		return chatResponse, status, err
	}

	thinkingProcess := utils.ExtractThinkingProcess(chatResponse.Message)
	_, err = s.Queries.CreateConversation(ctx, db.CreateConversationParams{
		SessionID: session.ID,
		Message:   thinkingProcess,
		Role:      "assistant",
	})
	chatResponse.Message = thinkingProcess
	if err != nil {
		return chatResponse, http.StatusInternalServerError, err
	}

	return chatResponse, status, nil
}

func (s *ChatService) isProductRelated(ctx context.Context, message string) bool {
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

func (s *ChatService) checkStockByKeyword(ctx context.Context, userMsg string) (string, error) {
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

func (s *ChatService) ChatCompletion(ctx context.Context, messages []model.ChatMessage) (model.ChatResponse, int, error) {
	return s.OpenAIService.ChatCompletion(ctx, messages)
}

func (s *ChatService) GetAllProductsAsString(ctx context.Context) (string, error) {
	products, err := s.Queries.GetAllProducts(ctx)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, p := range products {
		log.Println("Product value:", p.Name, p.Stock, p.Price)

		builder.WriteString(fmt.Sprintf("- %s (stok: %d, harga: Rp%d)\n", p.Name, p.Stock, p.Price))
	}

	return builder.String(), nil
}
