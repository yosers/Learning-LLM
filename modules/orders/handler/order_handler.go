package handler

import (
	"fmt"
	"net/http"
	order_model "shofy/modules/orders/model"
	"shofy/modules/orders/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) InitRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")	
	{
		orders.POST("/", h.CreateOrder)
		orders.GET("/", h.GetOrders)
		orders.GET("/:id", h.GetOrderById)
		orders.PUT("/:id", h.UpdateOrder)
		orders.DELETE("/:id", h.DeleteOrder)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}



