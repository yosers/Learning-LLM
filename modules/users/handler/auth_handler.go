package handler

import (
	"net/http"
	"shofy/modules/users/service"
	"shofy/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) InitRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/otp/send/:phone", h.SendOTP)
		authRoutes.POST("/otp/verify/:phone", h.VerifyOTP)
	}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	phoneNumber := c.Param("phone")

	// Validate phone number
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	// Generate and send OTP
	_, err := h.authService.GenerateAndSendOTP(c.Request.Context(), phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "OTP sent successfully", phoneNumber)

}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	phoneNumber := c.Param("phone")
	var input struct {
		OTP string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Success(c, http.StatusOK, "OTP is required", phoneNumber)
		return
	}

	// Verify OTP
	isValid, err := h.authService.VerifyOTP(c.Request.Context(), phoneNumber, input.OTP)

	if err != nil {
		response.Success(c, http.StatusOK, "System OTP Error", phoneNumber)
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", isValid)

}
