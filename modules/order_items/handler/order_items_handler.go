package handler

import (
	"net/http"
	middleware "shofy/middleware"
	"shofy/modules/order_items/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
)

type OrderItemsHandler struct {
	orderItemsService service.OrderItemsService
}

func NewOrderItemsHandler(orderItemsService service.OrderItemsService) *OrderItemsHandler {
	return &OrderItemsHandler{
		orderItemsService: orderItemsService,
	}
}

func (h *OrderItemsHandler) InitRoutes(router *gin.RouterGroup) {
	router.POST("/", middleware.RequireRole("ORDER_ITEMS_CREATE"), h.CreateOrderItems)
	router.GET("/:id", middleware.RequireRole("ORDER_ITEMS_GETBYID"), h.GetOrderItemsByID)
}

func (h *OrderItemsHandler) CreateOrderItems(c *gin.Context) {
	var req service.CreateOrderItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}
}
