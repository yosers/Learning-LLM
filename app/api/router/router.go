package router

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	chatHandler "shofy/modules/chat/handler"
	productHandler "shofy/modules/product/handler"
	pdService "shofy/modules/product/service"
	usHandler "shofy/modules/users/handler"
	usService "shofy/modules/users/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, srv *server.Server) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthCheck)

	v1Router := router.Group("/v1")

	chatRouter := chatHandler.NewChatAPIRoutes(ctx, srv)
	chatRouter.InitRoutes(v1Router)

	productService := pdService.NewProductService(srv.DBPool)
	handler := productHandler.NewProductHandler(productService)
	handler.InitRoutes(v1Router)

	userService := usService.NewUserService(srv.DBPool)
	userHandler := usHandler.NewUserHandler(userService)
	userHandler.InitRoutes(v1Router)

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
