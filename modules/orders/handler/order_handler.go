package handler

import (
	"fmt"
	"net/http"
	order_model "shofy/modules/orders/model"
	"shofy/modules/orders/service"
	"shofy/utils/response"
	"strconv"

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

	router.POST("/", h.CreateOrder)
	router.GET("/list", h.GetOrdersList)
	router.GET("/:id", h.GetOrderById)
	router.PUT("/:id", h.UpdateOrder)
	router.DELETE("/:id", h.DeleteOrder)

}

func (h *OrderHandler) GetOrdersList(c *gin.Context) {

	var q order_model.OrderQuery

	// Bind query parameters
	if err := c.BindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if q.Limit == 0 {
		q.Limit = 10
	}
	if q.CurrentPage == 0 {
		q.CurrentPage = 1
	}

	fmt.Println("q.UserID", q.UserID)
	fmt.Println("q.Status", q.Status)

	offset := (q.CurrentPage - 1) * q.Limit

	result, err := h.orderService.GetOrdersList(c.Request.Context(), int32(q.Limit), int32(offset), q.CurrentPage, int32(q.UserID), q.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Orders fetched successfully", gin.H{
		"data":         result.Items,
		"total_items":  result.TotalItems,
		"total_pages":  result.TotalPages,
		"current_page": result.CurrentPage,
		"limit":        result.Limit,
	})
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	_, err = h.orderService.GetOrderById(c.Request.Context(), int32(id))
	if err != nil {
		fmt.Println("Error: ", err)
		response.Error(c, http.StatusNotFound, "Order not found")
		return
	}

	order, err := h.orderService.GetOrderById(c.Request.Context(), int32(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get order")
		return
	}

	response.Success(c, http.StatusOK, "Order fetched successfully", order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	response.Success(c, http.StatusCreated, "Order created successfully", order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	_, err = h.orderService.GetOrderById(c.Request.Context(), int32(id))
	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Order not found", nil)
		return
	}

	var req service.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.ID = int32(id)

	order, err := h.orderService.UpdateOrder(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update order: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Order updated successfully", order)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	_, err = h.orderService.GetOrderById(c.Request.Context(), int32(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Order not found")
		return
	}

	err = h.orderService.DeleteOrder(c.Request.Context(), int32(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete order")
		return
	}
	response.Success(c, http.StatusOK, "Order deleted successfully", nil)
}
