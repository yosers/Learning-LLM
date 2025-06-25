package service

import (
	"context"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
)

type ChatServiceInterface interface {
	CreateChat(ctx context.Context, chat model.ChatPayload, channelID int) (model.ChatResponse, int, error)
	GetOrCreateSession(ctx context.Context, chatSession model.ChatSession) (db.Session, error)
	BuildMessageHistory(ctx context.Context, sessionID int32) ([]model.ChatMessage, error)
	SaveUserMessage(ctx context.Context, sessionID int32, content string) error
	SaveAssistantMessage(ctx context.Context, sessionID int32, content string) error
	ChatCompletion(ctx context.Context, messages []model.ChatMessage) (model.ChatResponse, int, error)
	GetAllProductsAsString(ctx context.Context) (string, error)
}
