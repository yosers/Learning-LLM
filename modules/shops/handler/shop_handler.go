package handler

import (
	"log"
	"net/http"
	model "shofy/modules/shops/model"
	"shofy/modules/shops/service"
	"shofy/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopService service.ShopService
}

func NewShopsHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

func (h *ShopHandler) InitRoutes(router *gin.RouterGroup) {
	router.GET("/", h.ListShops)
	router.POST("/", h.CreateShops)
	router.GET("/:id", h.GetShopsByID)
	router.PUT("/:id", h.UpdateShops)
	router.DELETE("/:id", h.DeleteShopsByID)
}

func (h *ShopHandler) ListShops(c *gin.Context) {
	var req model.ListShopRequest

	// Get query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = int32(page)

	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	req.PageSize = int32(pageSize)

	shopss, err := h.shopService.ListShops(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list shops:")
		return
	}

	response.Success(c, http.StatusOK, "Shops retrieved successfully", shopss)
}

func (h *ShopHandler) GetShopsByID(c *gin.Context) {
	shopsId := c.Param("id")

	shopsIdInt, err := strconv.Atoi(shopsId)
	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Invalid shops ID", nil)
		return
	}

	shops, err := h.shopService.GetShopsByID(c.Request.Context(), int32(shopsIdInt))
	if err != nil {
		if err.Error() == "shops not found" {
			response.NotSuccess(c, http.StatusOK, "Shops not found", nil)
			return
		}
		log.Printf("Error getting shops by ID %d: %v", shopsIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed to GetShopsByID")
		return
	}

	response.Success(c, http.StatusOK, "Shops retrieved successfully", shops)
}

func (h *ShopHandler) UpdateShops(c *gin.Context) {
	shopsId := c.Param("id")

	if shopsId == "" {
		response.Error(c, http.StatusBadRequest, "Shops ID is required")
		return
	}

	shopsIdInt, err := strconv.Atoi(shopsId)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Shops ID")
		return
	}

	var req model.ShopsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	shops, err := h.shopService.UpdateShops(c.Request.Context(), int32(shopsIdInt), &req)
	if err != nil {
		if err.Error() == "shops not found" {
			response.Error(c, http.StatusNotFound, "shops not found")
			return
		}
		log.Printf("Error UpdateShops shops by ID %d: %v", shopsIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed to update shops")
		return
	}

	response.Success(c, http.StatusOK, "Shops updated successfully", shops)
}

func (h *ShopHandler) DeleteShopsByID(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		response.NotSuccess(c, http.StatusOK, "Shops ID is required", nil)
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Invalid Shops ID", nil)
		return
	}

	if err := h.shopService.DeleteShopsByID(c.Request.Context(), int32(userIdInt)); err != nil {
		if err.Error() == "Shops not found" {
			response.NotSuccess(c, http.StatusOK, "Shops not found", nil)
			return
		}
		log.Printf("Error Delete Shops by ID %d: %v", userIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed DeleteShopsByID")
		return
	}

	response.Success(c, http.StatusOK, "Shops deleted successfully", nil)
}

func (h *ShopHandler) CreateShops(c *gin.Context) {
	var req model.ShopsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.shopService.CreateShops(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "Phone already exist" {
			response.Error(c, http.StatusConflict, "Phone already exist")
			return
		}
		log.Printf("Error CreateShops", err)
		response.Error(c, http.StatusInternalServerError, "Failed to create Shops")
		return
	}

	response.Success(c, http.StatusCreated, "Shops created successfully", user)
}
