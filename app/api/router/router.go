package router

import (
	"context"
	"net/http"
	"shofy/app/api/server"

	handler "shofy/modules/chat/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, srv *server.Server) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthCheck)

	v1Router := router.Group("/v1")

	chatRouter := handler.NewChatAPIRoutes(ctx, srv)
	chatRouter.InitRoutes(v1Router)
	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
