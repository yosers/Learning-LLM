package handler

import (
	"fmt"
	"net/http"
	"shofy/modules/role/service"
	"shofy/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) InitRoutes(router *gin.RouterGroup) {
	router.GET("/list", h.ListRoles)
	router.POST("/", h.CreateRoles)
	router.PUT("/", h.UpdateRolesById)
	router.GET("/:id", h.RolesByID)
	router.DELETE("/:id", h.DeleteRolesByID)
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleService.ListRole(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Success(c, http.StatusOK, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) RolesByID(c *gin.Context) {
	idRole, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleService.RolesByID(c.Request.Context(), int32(idRole))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if role == nil {
		response.Error(c, http.StatusNotFound, "Role not found")
		return
	}

	response.Success(c, http.StatusOK, "Role retrieved successfully", role)
}

func (h *RoleHandler) DeleteRolesByID(c *gin.Context) {
	idRole, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.roleService.DeleteRolesByID(c.Request.Context(), int32(idRole))
	if err != nil {
		if err.Error() == "role not found" {
			response.Error(c, http.StatusNotFound, "Role not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete role: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Role deleted successfully", nil)
}

func (h *RoleHandler) UpdateRolesById(c *gin.Context) {
	var roleUpdate service.UpdateRolesRequest

	if err := c.ShouldBindJSON(&roleUpdate); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.roleService.UpdateRolesById(c.Request.Context(), &roleUpdate)

	if err != nil {
		if err.Error() == "role not found" {
			response.Error(c, http.StatusNotFound, "Role not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Role updated successfully", nil)
}

func (h *RoleHandler) CreateRoles(c *gin.Context) {
	var roleCreate service.CreateRolesRequest

	if err := c.ShouldBindJSON(&roleCreate); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, err := h.roleService.CreateRoles(c.Request.Context(), &roleCreate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create role: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Role created successfully", role)
}
