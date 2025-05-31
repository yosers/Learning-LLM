package handler

import (
	"net/http"
	"shofy/modules/categories/service"
	"shofy/utils/response"

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
	categories := router.Group("/categories")
	{
		categories.GET("", h.GetAllCategories)
		categories.DELETE("/:id", h.DeleteCategory)
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

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {

	response.Success(c, http.StatusOK, "Successfully retrieved categories", categories)
}
