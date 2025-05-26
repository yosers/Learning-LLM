package service

import (
	"context"
	"shofy/modules/chat/model"
)

type ChatServiceInterface interface {
	CreateChat(ctx context.Context, chat model.ChatPayload, channelID int) (model.ChatResponse, int, error)
}
