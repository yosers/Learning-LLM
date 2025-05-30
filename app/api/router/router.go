package router

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	"shofy/middleware"
	chatHandler "shofy/modules/chat/handler"
	productHandler "shofy/modules/product/handler"
	pdService "shofy/modules/product/service"
	usHandler "shofy/modules/users/handler"
	usService "shofy/modules/users/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, srv *server.Server) *gin.Engine {
	router := gin.Default()

	// Public endpoints
	router.GET("/health", healthCheck)

	v1Router := router.Group("/v1")

	// Public routes
	userService := usService.NewUserService(srv.DBPool)
	userHandler := usHandler.NewUserHandler(userService)
	userHandler.InitRoutes(v1Router)

	// Chat routes
	chatRouter := chatHandler.NewChatAPIRoutes(ctx, srv)
	chatRouter.InitRoutes(v1Router)

	// Product routes (tanpa autentikasi)
	// productService := pdService.NewProductService(srv.DBPool)
	// productHandler := productHandler.NewProductHandler(productService)
	// productHandler.InitRoutes(v1Router.Group("/products"))

	// Protected routes
	protectedRoutes := v1Router.Group("")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		// Product routes
		productService := pdService.NewProductService(srv.DBPool)
		productHandler := productHandler.NewProductHandler(productService)
		productHandler.InitRoutes(protectedRoutes.Group("/products"))
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
