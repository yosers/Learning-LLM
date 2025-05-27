package router

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	chatHandler "shofy/modules/chat/handler"
	productHandler "shofy/modules/product/handler"
	"shofy/modules/product/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, srv *server.Server) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthCheck)

	v1Router := router.Group("/v1")

	chatRouter := chatHandler.NewChatAPIRoutes(ctx, srv)
	chatRouter.InitRoutes(v1Router)

	productService := service.NewProductService(srv.DBPool)
	handler := productHandler.NewProductHandler(productService)
	handler.InitRoutes(v1Router)

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
