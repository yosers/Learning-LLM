package handler

import (
	"net/http"
	"shofy/middleware"
	"shofy/modules/categories/service"
	"shofy/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) InitRoutes(router *gin.RouterGroup) {

	categories := router.Group("/")

	protected := categories.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		router.GET("/list/limit/:limit/offset/:offset", h.GetCategoriesPaginated)
		router.GET("/all", h.GetAllCategories)
		router.GET("/detail/:id", h.GetCategoryByID)
		router.DELETE("/delete/:id", h.DeleteCategoryByID)
	}
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategory(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Successfully retrieved categories", categories)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}
	categories, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", gin.H{
		"category": categories,
	})
}

func (h *CategoryHandler) GetCategoriesPaginated(c *gin.Context) {
	lim := c.Param("limit")
	off := c.Param("offset")

	limitStr := c.DefaultQuery("limit", lim)
	offsetStr := c.DefaultQuery("offset", off)

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	categories, err := h.categoryService.GetCategoriesPaginated(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Categories retrieved successfully", gin.H{
		"data":         categories.Items,
		"total_items":  categories.TotalItems,
		"total_pages":  categories.TotalPages,
		"current_page": categories.CurrentPage,
		"limit":        categories.Limit,
	})
}

func (h *CategoryHandler) DeleteCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err := h.categoryService.DeleteByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Category deleted successfully", nil)
}

// func (h *CategoryHandler) DeleteCategory(c *gin.Context) {

// 	response.Success(c, http.StatusOK, "Successfully retrieved categories", categories)
// }
