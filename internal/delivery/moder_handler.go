package delivery

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"userService/internal/model"
	"userService/internal/service"
	"userService/pkg/error/format"
)

type ModerHandler struct {
	ms *service.ModerService
}

func NewModerHandler(ms *service.ModerService) *ModerHandler {
	return &ModerHandler{ms: ms}
}

// SetRoleById assigns a new role to a user.
// @Tags moder
// @Router /api/v1/moder/role/{id} [put]
func (h *ModerHandler) SetRoleById(c *gin.Context) {
	editorID, editorRole, userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	var req struct {
		RoleName string `json:"role_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	newRole := model.Role(req.RoleName)
	if err := h.ms.SetRoleById(editorID, editorRole, userID, newRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User role updated successfully"})
}

// BanUserById bans a user.
// @Tags moder
// @Router /api/v1/moder/ban/{id} [put]
func (h *ModerHandler) BanUserById(c *gin.Context) {
	_, editorRole, userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	if err := h.ms.BanUserById(editorRole, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User banned successfully"})
}

// UnBanUserById unbans a user.
// @Tags moder
// @Router /api/v1/moder/unban/{id} [put]
func (h *ModerHandler) UnBanUserById(c *gin.Context) {
	_, editorRole, userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	if err := h.ms.UnBanUserById(editorRole, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User unbanned successfully"})
}

// validateModeration extracts user ID and role from the context.
func validateModeration(c *gin.Context) (int, model.Role, int, error) {
	editorIDRaw, exists := c.Get("user_id")
	if !exists {
		return 0, "", 0, errors.New("user is not logged in")
	}

	editorID, ok := editorIDRaw.(int)
	if !ok {
		return 0, "", 0, errors.New("invalid editor ID format")
	}

	editorRoleRaw, exists := c.Get("user_role")
	if !exists {
		return 0, "", 0, errors.New("user role not found")
	}

	editorRole, ok := editorRoleRaw.(model.Role)
	if !ok {
		return 0, "", 0, errors.New("invalid user role format")
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, "", 0, errors.New("invalid user ID")
	}

	if editorID == userID {
		return 0, "", 0, errors.New("cannot moderate own account")
	}

	return editorID, editorRole, userID, nil
}
