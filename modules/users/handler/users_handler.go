package handler

import (
	"net/http"
	"shofy/modules/users/service"
	"shofy/utils/response"
	"strconv"

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
		users.GET("/verify-otp/:otp/:user-id", h.VerifyOTP)

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
		response.Error(c, http.StatusOK, "Failed check phone number")
		return
	}

	response.Success(c, http.StatusOK, "Check Phone successfully", users)
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	otpUser := c.Param("otp")
	userId := c.Param("user-id")

	if otpUser == "" {
		response.Error(c, http.StatusOK, "OTP is required")
		return
	}

	if userId == "" {
		response.Error(c, http.StatusOK, "User id is required")
		return
	}

	userIdNumber, err := strconv.Atoi(userId)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	result, err := h.userService.VerifyOTP(c.Request.Context(), otpUser, userIdNumber)

	if err != nil {
		response.Error(c, http.StatusOK, "Failed to Verify OTP")
		return
	}

	// Set token in cookie
	c.SetCookie("token", result.Token, 24*60*60, "/", "", false, true) // 24 hours expiry

	response.Success(c, http.StatusOK, "Verify OTP successfully", result)
}
