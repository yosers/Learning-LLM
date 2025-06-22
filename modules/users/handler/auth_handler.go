package handler

import (
	"log"
	"net/http"
	"shofy/modules/users/service"
	"shofy/utils/response"
	"strings"

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
		response.NotSuccess(c, http.StatusBadRequest, "Phone number is required", nil)
		return
	}

	if req.Code == "" {
		response.NotSuccess(c, http.StatusBadRequest, "Code is required", nil)
		return
	}

	// Generate and send OTP
	data, err := h.authService.GenerateAndSendOTP(c.Request.Context(), req)
	if err != nil {
		// Check if the error is "Data Not Found"
		if strings.Contains(err.Error(), "Data Not Found") {
			response.Error(c, http.StatusNotFound, "Data not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if data.Remarks == "User already logged" {
		response.NotSuccess(c, http.StatusOK, "User already logged", nil)
		return
	}

	response.Success(c, http.StatusOK, "OTP sent successfully", nil)

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
		if strings.Contains(err.Error(), "OTP conflict: multiple valid entries found") {
			response.NotSuccess(c, http.StatusOK, "Multiple OTP", "OTP conflict: Please request a new OTP")
			return
		} else if strings.Contains(err.Error(), "OTP has expired") {
			response.NotSuccess(c, http.StatusOK, "OTP Expired", "Please request a new OTP")
			return
		} else if strings.Contains(err.Error(), "OTP has already been used") {
			response.NotSuccess(c, http.StatusOK, "User already logged", "User already logged")
			return
		}
		response.NotSuccess(c, http.StatusOK, "System OTP Error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", isValid)

}
