package handler

import (
	"net/http"
	"shofy/modules/users/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) InitRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/phone/:phone", h.GetUserPhoneNumber)
	}
}

func (h *UserHandler) GetUserPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phone")

	if phoneNumber == "" {
		response.Error(c, http.StatusOK, "Phone number is required")
		return
	}

	// Convert phoneNumber string to pgtype.Text
	phoneNumberText := pgtype.Text{String: phoneNumber, Valid: true}
	users, err := h.userService.GetUserPhoneNumber(c.Request.Context(), phoneNumberText)
	if err != nil {
		response.Error(c, http.StatusOK, "Failed to check phone number")
		return
	}

	response.Success(c, http.StatusOK, "Check Phone successfully", users)
}
