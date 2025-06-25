package handler

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
	chatService "shofy/modules/chat/service"
	deepinfraService "shofy/modules/deepinfra/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRouter struct {
	Query         *db.Queries
	DBPool        *pgxpool.Pool
	ChatService   chatService.ChatServiceInterface
	OpenAIService *deepinfraService.OpenAIService
}

func NewChatAPIRoutes(ctx context.Context, srv *server.Server) *ChatRouter {
	chatSvc := chatService.NewChatService(ctx, srv.DBPool, srv.Queries)
	return &ChatRouter{
		Query:       srv.Queries,
		DBPool:      srv.DBPool,
		ChatService: chatSvc,
	}
}

func (r *ChatRouter) InitRoutes(rg *gin.RouterGroup) {
	rg.POST("/chat", r.CreateChat)
	rg.POST("/chat/session", r.GetOrCreateSession)
	rg.POST("/chat/message", r.MessageChat)
	// rg.GET("/chat/session/messages", r.CreateChat)

}

func (r *ChatRouter) CreateChat(c *gin.Context) {
	ctx := c.Request.Context()
	var chatPayload model.ChatPayload

	if err := c.ShouldBindJSON(&chatPayload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, _, err := r.ChatService.CreateChat(ctx, chatPayload, chatPayload.ChannelID)

	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Internal Server Error", nil)
		return
	}

	response.Success(c, http.StatusOK, "Chat Successfully", result.Message)
}

func (r *ChatRouter) GetOrCreateSession(c *gin.Context) {
	ctx := c.Request.Context()
	var chatPayload model.ChatSession

	if err := c.ShouldBindJSON(&chatPayload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := r.ChatService.GetOrCreateSession(ctx, chatPayload)

	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Internal Server Error", nil)
		return
	}

	response.Success(c, http.StatusOK, "Chat Successfully", result)
}

func (r *ChatRouter) MessageChat(c *gin.Context) {
	ctx := c.Request.Context()
	var payload model.ChatMessagePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Ambil histori chat
	history, err := r.ChatService.BuildMessageHistory(ctx, payload.SessionID)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal mengambil histori", nil)
		return
	}

	// Simpan pesan user ke DB
	err = r.ChatService.SaveUserMessage(ctx, payload.SessionID, payload.Message)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal menyimpan pesan user", nil)
		return
	}

	// Tambahkan pesan user terbaru ke history
	history = append(history, model.ChatMessage{
		Role:    "user",
		Content: payload.Message,
	})

	// 💡 Panggil GetAllProductsAsString DI SINI
	productsStr, err := r.ChatService.GetAllProductsAsString(ctx)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal mengambil data produk", nil)
		return
	}

	// 💡 Sisipkan prompt system tentang produk
	systemPrompt := model.ChatMessage{
		Role:    "system",
		Content: "Berikut daftar produk tersedia:\n" + productsStr,
	}

	// Tambahkan system prompt ke awal
	history = append([]model.ChatMessage{systemPrompt}, history...)

	// Kirim ke AI
	reply, _, err := r.ChatService.ChatCompletion(ctx, history)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal mendapatkan jawaban dari AI", nil)
		return
	}

	// Simpan jawaban AI
	err = r.ChatService.SaveAssistantMessage(ctx, payload.SessionID, reply.Message)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal menyimpan jawaban AI", nil)
		return
	}

	// Kirim ke client
	response.Success(c, http.StatusOK, "Berhasil membalas pesan", reply.Message)
}
