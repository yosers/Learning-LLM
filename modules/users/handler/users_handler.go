package handler

import (
	"net/http"
	"shofy/middleware"
	"shofy/modules/users/service"
	"shofy/utils/response"
	"strconv"
	"strings"

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

		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", h.Logout)
		}
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

	users, err := h.userService.GenerateOTPByPhone(c.Request.Context(), phoneNumberText)

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

func (h *UserHandler) Logout(c *gin.Context) {
	// Get token from cookie or Authorization header
	token := ""

	// Try to get from cookie first
	if cookieToken, err := c.Cookie("token"); err == nil && cookieToken != "" {
		token = cookieToken
	} else {
		// If no cookie, try Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		response.Error(c, http.StatusBadRequest, "No token provided")
		return
	}

	// Invalidate the token
	if err := h.userService.Logout(c.Request.Context(), token); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Clear the cookie if it exists
	c.SetCookie("token", "", -1, "/", "", false, true)

	response.Success(c, http.StatusOK, "Successfully logged out", nil)
}
