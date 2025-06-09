package router

import (
	"context"
	"net/http"
	"shofy/app/api/server"
	categoryHandler "shofy/modules/categories/handler"
	categoryService "shofy/modules/categories/service"
	chatHandler "shofy/modules/chat/handler"
	orderHandler "shofy/modules/orders/handler"
	orderService "shofy/modules/orders/service"
	productHandler "shofy/modules/product/handler"
	pdService "shofy/modules/product/service"
	usHandler "shofy/modules/users/handler"
	usService "shofy/modules/users/service"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx context.Context, srv *server.Server) *gin.Engine {
	router := gin.Default()

	// ✅ Tambahkan middleware CORS DI SINI
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // FE and BE addresses
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public endpoints
	router.GET("/health", healthCheck)

	v1Router := router.Group("/v1")

	userService := usService.NewUserService(srv.DBPool)
	userHandler := usHandler.NewUserHandler(userService)
	userHandler.InitRoutes(v1Router)

	// Public routes
	authService := usService.NewAuthService(srv.DBPool)
	authHandler := usHandler.NewAuthHandler(authService)
	authHandler.InitRoutes(v1Router)

	// Chat routes
	chatRouter := chatHandler.NewChatAPIRoutes(ctx, srv)
	chatRouter.InitRoutes(v1Router)

	// Product routes (tanpa autentikasi)
	productService := pdService.NewProductService(srv.DBPool)
	productHandler := productHandler.NewProductHandler(productService)
	productHandler.InitRoutes(v1Router.Group("/products"))

	categoryService := categoryService.NewCategoryService(srv.DBPool)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)
	categoryHandler.InitRoutes(v1Router.Group("/categories"))

	orderService := orderService.NewOrderService(srv.DBPool)
	orderHandler := orderHandler.NewOrderHandler(orderService)
	orderHandler.InitRoutes(v1Router.Group("/orders"))

	// Protected routes
	// protectedRoutes := v1Router.Group("")
	// protectedRoutes.Use(middleware.AuthMiddleware())
	// {
	// 	// Product routes
	// 	productService := pdService.NewProductService(srv.DBPool)
	// 	productHandler := productHandler.NewProductHandler(productService)
	// 	productHandler.InitRoutes(protectedRoutes.Group("/products"))
	// }

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
