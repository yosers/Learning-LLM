package handler

import (
	"net/http"
	"strconv"

	middleware "shofy/middleware"
	"shofy/modules/product/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) InitRoutes(router *gin.RouterGroup) {
	// All these routes will inherit the middleware from the group
	router.GET("/list", middleware.RequireRole("PRODUCT_LIST"), h.ListProducts) // Changed from "" to "/list" for clarity
	router.GET("/all", h.GetAllProducts)
	router.GET("/detail/:id", middleware.RequireRole("PRODUCT_VIEW"), h.GetProductByID)       // Changed from ":id" to "/detail/:id" for clarity
	router.DELETE("/deleted/:id", middleware.RequireRole("PRODUCT_DELETE"), h.GetProductByID) // Changed from ":id" to "/detail/:id" for clarity

}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", gin.H{
		"user_id": userID,
		"product": product,
	})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	response.Error(c, http.StatusUnauthorized, "User not authenticated")
	// 	return
	// }

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * limit

	result, err := h.productService.ListProducts(c.Request.Context(), int32(limit), int32(offset), page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", gin.H{
		// "user_id":      userID,
		"data":         result.Items,
		"total_items":  result.TotalItems,
		"total_pages":  result.TotalPages,
		"current_page": result.CurrentPage,
		"limit":        result.Limit,
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	response.Error(c, http.StatusUnauthorized, "User not authenticated")
	// 	return
	// }

	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "All products retrieved successfully", gin.H{
		// "user_id":  userID,
		"products": products,
	})
}
