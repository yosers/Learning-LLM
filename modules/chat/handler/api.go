package handler

import (
	"context"
	"fmt"
	"net/http"
	"shofy/app/api/server"
	db "shofy/db/sqlc"
	"shofy/modules/chat/model"
	chatService "shofy/modules/chat/service"
	deepinfraService "shofy/modules/deepinfra/service"
	"shofy/utils/response"
	"strings"

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

	shotpStr, err := r.ChatService.GetAllProductsAsString(ctx)
	if err != nil {
		response.NotSuccess(c, http.StatusInternalServerError, "Gagal mengambil data produk", nil)
		return
	}

	// 💡 Sisipkan prompt system tentang produk
	// System Prompt — menjelaskan role AI dan produk
	systemPrompt := model.ChatMessage{
		Role: "system",
		Content: fmt.Sprintf(`Kamu adalah asisten virtual dari sebuah toko online bernama "Shofy".
					Tugas kamu:
					- Menjawab pertanyaan tentang produk yang dijual
					- Menjelaskan detail dan stok barang
					- Membantu pelanggan dalam proses pemesanan
					- Memberikan jawaban yang relevan dan informatif

					Berikan jawaban yang rapi dengan format seperti berikut:

					Produk "sepatu" tersedia di:

					1. Toko A
					- Stok: 10
					- Harga: Rp 100.000

					2. Toko B
					- Stok: 20
					- Harga: Rp 120.000

					Pisahkan setiap item dengan newline (\n) agar mudah dibaca di frontend.

					Data produk:
					%s

					Data toko:
					%s`, productsStr, shotpStr),
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

	// Ubah newline menjadi <br> sebelum dikirim ke frontend
	formattedReply := strings.ReplaceAll(reply.Message, "\n", "<br>")

	// Kirim ke client
	response.Success(c, http.StatusOK, "Berhasil membalas pesan", formattedReply)
}
