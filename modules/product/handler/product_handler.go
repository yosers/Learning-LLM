package handler

import (
	"fmt"
	"net/http"
	product_model "shofy/modules/product/model"
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
	router.GET("/", h.ListProducts) // Changed from "" to "/list" for clarity
	//router.GET("/all", h.GetAllProducts)
	router.POST("/", h.CreateProduct)
	router.GET("/:id", h.GetProductByID)

	//router.PUT("/:id", h.UpdateProduct)                                                  // Changed from ":id" to "/detail/:id" for clarity
	//router.DELETE("/:id", middleware.RequireRole("PRODUCT_DELETE"), h.DeleteProductByID) // Changed from ":id" to "/detail/:id" for clarity

	router.PUT("/:id", h.UpdateProduct)        // Changed from ":id" to "/detail/:id" for clarity
	router.DELETE("/:id", h.DeleteProductByID) // Changed from ":id" to "/detail/:id" for clarity

}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create product: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.ID = id

	product, err := h.productService.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update product: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)

	// All these routes will inherit the middleware from the group

}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	//  exists = c.Get("user_id")
	// if !exists {
	// 	response.Error(c, http.StatusUnauthorized, "User not authenticated")
	// 	return
	// }

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
		//"user_id": userID,
		"product": product,
	})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var q product_model.ProductQuery

	// Bind query parameters
	if err := c.BindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Set defaults if not provided
	if q.Limit == 0 {
		q.Limit = 10
	}
	if q.CurrentPage == 0 {
		q.CurrentPage = 1
	}

	offset := (q.CurrentPage - 1) * q.Limit

	result, err := h.productService.ListProducts(c.Request.Context(), int32(q.Limit), int32(offset), q.CurrentPage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", gin.H{
		"product":      result.Items,
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
		//"user_id":  userID,
		"products": products,
	})
}

func (h *ProductHandler) DeleteProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Check if product exists
	_, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	// Proceed to delete
	err = h.productService.DeleteProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}
