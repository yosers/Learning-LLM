package handler

import (
	"log"
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
		authRoutes.POST("/otp/send", h.SendOTP)
		authRoutes.POST("/otp/verify", h.VerifyOTP)
	}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	log.Println("WHATSHAP")

	var req service.SendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate phone number
	if req.Phone == "" {
		response.Success(c, http.StatusBadRequest, "Phone number is required", req.Phone)
		return
	}

	if req.Code == "" {
		response.Success(c, http.StatusBadRequest, "Code is required", req.Phone)
		return
	}

	// Generate and send OTP
	_, err := h.authService.GenerateAndSendOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, http.StatusOK, "OTP sent successfully", req.Phone)

}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {

	var input service.VerifyOTP

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify OTP
	isValid, err := h.authService.VerifyOTP(c.Request.Context(), input.Otp)

	if err != nil {
		response.Success(c, http.StatusOK, "System OTP Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", isValid)

}
