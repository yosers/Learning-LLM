package service

import (
	"context"
	"net/http"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"

	azureService "shofy/modules/azure/service"

	openaiService "shofy/modules/gpt/service"

	utils "shofy/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openai/openai-go"
)

type ChatService struct {
	DBPool  *pgxpool.Pool
	Queries *db.Queries

	// GPTService *service.GPTService
	// AzureOpenAI *azureService.AzureOpenAI
	OpenAIService *openaiService.OpenAIService
}

func NewChatService(ctx context.Context, dbPool *pgxpool.Pool, http *http.Client, queries *db.Queries, azureOpenAI *azureService.AzureOpenAI) *ChatService {
	// gptService := service.NewGPTService(ctx, dbPool, http, queries)
	openaiService := openaiService.NewOpenAIService(ctx)
	return &ChatService{
		DBPool:        dbPool,
		Queries:       queries,
		OpenAIService: openaiService,
		// AzureOpenAI: azureOpenAI,
		// GPTService: gptService,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, chat model.ChatPayload, channelID int) (model.ChatResponse, int, error) {
	session, err := s.Queries.GetCurrentSessions(ctx, db.GetCurrentSessionsParams{
		ChannelID: int32(channelID),
		UserID:    1,
	})
	//
	var hasNoSession bool = false
	messages := []openai.ChatCompletionMessageParamUnion{}
	if err != nil {
		// no session found, create a new session
		session, err = s.Queries.CreateSession(ctx, db.CreateSessionParams{
			ChannelID: int32(channelID),
			UserID:    1,
		})
		if err != nil {
			return model.ChatResponse{}, http.StatusInternalServerError, err
		}
		hasNoSession = true
	} else {
		// Get conversation by session id
		conversations, err := s.Queries.GetConversationsBySessionID(ctx, session.ID)
		if err != nil {
			// do nothing if no conversation found
			hasNoSession = true
		}

		for _, conversation := range conversations {
			if conversation.Role == "user" {
				messages = append(messages, openai.UserMessage(conversation.Message))
			} else if conversation.Role == "assistant" {
				messages = append(messages, openai.AssistantMessage(conversation.Message))
			}
		}
	}

	conversation, err := s.Queries.CreateConversation(ctx, db.CreateConversationParams{
		SessionID: session.ID,
		Message:   chat.Message,
		Role:      "user",
	})

	if hasNoSession {
		messages = append(messages, openai.SystemMessage(utils.PromptFirstMessage))
	}

	messages = append(messages, openai.UserMessage(conversation.Message))

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
