package handler

import (
	"log"
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
	router.GET("/", h.ListRoles)
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
		if err.Error() == "Role not found" {
			response.NotSuccess(c, http.StatusOK, "Role not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
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
			response.NotSuccess(c, http.StatusOK, "Role not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete role")
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
			response.NotSuccess(c, http.StatusOK, "Role not found", nil)
			return
		}
		log.Printf("Error Failed to update Roles", err)
		response.Error(c, http.StatusInternalServerError, "Failed to update user")
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
		log.Printf("Error creating role: %v", err)
		response.Error(c, http.StatusInternalServerError, "Failed to create role")
		return
	}

	response.Success(c, http.StatusCreated, "Role created successfully", role)
}
