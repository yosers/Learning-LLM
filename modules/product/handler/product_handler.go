package handler

import (
	"net/http"
	"strconv"

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
	products := router.Group("/products")
	{
		products.GET("", h.ListProducts)
		products.GET("/all", h.GetAllProducts)
		products.GET("/:id", h.GetProductByID)
	}
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
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

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * limit

	result, err := h.productService.ListProducts(c.Request.Context(), int32(limit), int32(offset), page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", gin.H{
		"message":      "Products retrieved successfully",
		"data":         result.Items,
		"total_items":  result.TotalItems,
		"total_pages":  result.TotalPages,
		"current_page": result.CurrentPage,
		"limit":        result.Limit,
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "All products retrieved successfully", products)
}
