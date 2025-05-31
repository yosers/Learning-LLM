package handler

import (
	"fmt"
	"net/http"
	"shofy/middleware"
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
		users.GET("/phone/:phone", h.GenerateOTPByPhone)
		users.GET("/verify-otp/:otp/:user-id", h.VerifyOTP)

		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			users.POST("/save", h.CreateUser)
			users.GET("/list", h.ListUsers)
			protected.POST("/logout/:user-id", h.Logout)
			protected.PUT("/:user-id", h.UpdateUser)
		}
	}
}

func (h *UserHandler) GenerateOTPByPhone(c *gin.Context) {
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
	userId := c.Param("user-id")

	// Try to get from cookie first
	if cookieToken, err := c.Cookie("token"); err == nil && cookieToken != "" {
		token = cookieToken
	} else {
		// If no cookie, try Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	if token == "" {
		response.Error(c, http.StatusBadRequest, "No token provided")
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

	// Invalidate the token
	if err := h.userService.Logout(c.Request.Context(), token, userIdNumber); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Clear the cookie if it exists
	c.SetCookie("token", "", -1, "/", "", false, true)

	response.Success(c, http.StatusOK, "Successfully logged out", nil)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("user-id")
	if userId == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), int32(userIdInt), &req)
	if err != nil {
		if err.Error() == "user not found" {
			response.Error(c, http.StatusNotFound, "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var req service.ListUsersRequest

	// Get query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = int32(page)

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	req.PageSize = int32(pageSize)

	shopId, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid shop ID")
		return
	}
	req.ShopID = int32(shopId)

	users, err := h.userService.ListUsers(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to list users: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}
