package handler

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	db "shofy/db/sqlc"

	"shofy/modules/chat/model"
	chatService "shofy/modules/chat/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRouter struct {
	Query       *db.Queries
	DBPool      *pgxpool.Pool
	ChatService chatService.ChatServiceInterface
}

func NewChatAPIRoutes(ctx context.Context, srv *server.Server) *ChatRouter {
	chatService := chatService.NewChatService(ctx, srv.DBPool, nil, srv.Queries, srv.OpenAI)
	return &ChatRouter{
		Query:       srv.Queries,
		DBPool:      srv.DBPool,
		ChatService: chatService,
	}
}

func (r *ChatRouter) InitRoutes(rg *gin.RouterGroup) {
	rg.POST("/chat", r.CreateChat)
}

func (r *ChatRouter) CreateChat(c *gin.Context) {
	ctx := c.Request.Context()
	var chatPayload model.ChatPayload
	if err := c.ShouldBindJSON(&chatPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := r.ChatService.CreateChat(ctx, chatPayload, chatPayload.ChannelID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, map[string]interface{}{
		"message": "success",
		"data":    result.Message,
		"status":  status,
	})
}
