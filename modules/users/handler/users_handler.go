package handler

import (
	"log"
	"net/http"
	"shofy/modules/users/service"
	"shofy/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
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
	router.GET("/list", h.ListUsers)
	router.POST("/", h.CreateUser)
	router.POST("/logout", h.Logout)
	router.PUT("/:id", h.UpdateUser)
	router.GET("/:id", h.GetUsersByID)
	router.DELETE("/:id", h.DeleteUsersByID)

}

func (h *UserHandler) Logout(c *gin.Context) {
	var req service.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get token from cookie or Authorization header
	token := ""

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

	userIdNumber, err := strconv.Atoi(req.UserID)

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
		if err.Error() == "Phone already exist" {
			response.Error(c, http.StatusConflict, "Phone already exist")
			return
		}
		log.Printf("Error CreateUser", err)
		response.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("id")

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
		log.Printf("Error UpdateUser user by ID %d: %v", userIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed to update user")
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

	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	req.PageSize = int32(pageSize)

	// shopId, err := strconv.Atoi(c.Query("shop_id"))
	// if err != nil {
	// 	response.Error(c, http.StatusBadRequest, "Invalid shop ID")
	// 	return
	// }
	// req.ShopID = int32(shopId)

	users, err := h.userService.ListUsers(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users:")
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUsersByID(c *gin.Context) {
	userId := c.Param("id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), int32(userIdInt))
	if err != nil {
		if err.Error() == "user not found" {
			response.NotSuccess(c, http.StatusOK, "User not found", nil)
			return
		}
		log.Printf("Error getting user by ID %d: %v", userIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed to GetUsersByID")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) DeleteUsersByID(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		response.NotSuccess(c, http.StatusOK, "User ID is required", nil)
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		response.NotSuccess(c, http.StatusOK, "Invalid user ID", nil)
		return
	}

	if err := h.userService.DeleteUsersByID(c.Request.Context(), int32(userIdInt)); err != nil {
		if err.Error() == "user not found" {
			response.NotSuccess(c, http.StatusOK, "User not found", nil)
			return
		}
		log.Printf("Error Delete user by ID %d: %v", userIdInt, err)
		response.Error(c, http.StatusInternalServerError, "Failed DeleteUsersByID")
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
