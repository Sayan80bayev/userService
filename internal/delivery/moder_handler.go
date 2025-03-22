package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"userService/internal/model"
	"userService/internal/service"
)

type ModerHandler struct {
	ms *service.ModerService
}

func NewModerHandler(ms *service.ModerService) *ModerHandler {
	return &ModerHandler{ms: ms}
}

// SetRoleById sets the role for a user by ID.
// @Tags moder
// @Router /api/v1/moder/role/{id} [put]
func (h *ModerHandler) SetRoleById(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var req struct {
		RoleName string `json:"role_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	err = h.ms.SetRoleById(userID, model.Role(req.RoleName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to set role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User role updated successfully"})
}

// BanUserById bans a user by ID.
// @Tags moder
// @Router /api/v1/moder/ban/{id} [put]
func (h *ModerHandler) BanUserById(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	err = h.ms.BanUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to ban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User banned successfully"})
}

// UnBanUserById unbans a user by ID.
// @Tags moder
// @Router /api/v1/moder/unban/{id} [put]
func (h *ModerHandler) UnBanUserById(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	err = h.ms.UnBanUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to unban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User unbanned successfully"})
}
