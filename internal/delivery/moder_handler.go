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

// SetRoleById assigns a role to a user using their ID.
// @Tags moder
// @Router /api/v1/moder/role/{id} [put]
func (h *ModerHandler) SetRoleById(c *gin.Context) {
	userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	var req struct {
		RoleName string `json:"role_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body."})
		return
	}

	role := model.Role(req.RoleName)
	if !isValidRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid role name."})
		return
	}

	if err := h.ms.SetRoleById(userID, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to assign the role."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User role updated successfully."})
}

// BanUserById bans a user using their ID.
// @Tags moder
// @Router /api/v1/moder/ban/{id} [put]
func (h *ModerHandler) BanUserById(c *gin.Context) {
	userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	if err := h.ms.BanUserById(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to ban user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User banned successfully."})
}

// UnBanUserById unbans a user using their ID.
// @Tags moder
// @Router /api/v1/moder/unban/{id} [put]
func (h *ModerHandler) UnBanUserById(c *gin.Context) {
	userID, err := validateModeration(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": format.CapitalizeError(err)})
		return
	}

	if err := h.ms.UnBanUserById(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to unban user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User unbanned successfully."})
}

// isValidRole checks whether the given role is valid.
func isValidRole(role model.Role) bool {
	return role == model.RoleAdmin || role == model.RoleModerator || role == model.RoleUser
}

// validateModeration checks whether the current moderator has permission to perform this action.
func validateModeration(c *gin.Context) (int, error) {
	editorIDRaw, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user is not logged in")
	}

	editorID, ok := editorIDRaw.(int)
	if !ok {
		return 0, errors.New("invalid editor ID format")
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("invalid user ID")
	}

	if editorID == userID {
		return 0, errors.New("cannot moderate own account")
	}

	return userID, nil
}
