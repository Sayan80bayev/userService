package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"userService/internal/service"
	"userService/internal/transfer/request"
	"userService/pkg/logging"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

// UpdateUser обновляет информацию о пользователе
// @Summary Обновление пользователя
// @Description Позволяет обновить информацию о пользователе, включая аватар
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param userId header string true "ID пользователя"
// @Param avatar formData file false "Аватар пользователя"
// @Param user body request.UserRequest true "Данные пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [put]
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    "UNAUTHORIZED",
			"message": "You're unauthorized",
		})
		return
	}

	var ur request.UserRequest
	if err := ctx.ShouldBind(&ur); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_INPUT",
			"message": "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	avatar, header, err := ctx.Request.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_INPUT",
			"message": "Could not get avatar",
			"details": err.Error(),
		})
		logging.Instance.Warn("Error on getting avatar", err.Error())
		return
	}
	ur.Avatar, ur.Header = avatar, header

	err = h.service.UpdateUser(ur, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "SERVER_ERROR",
			"message": "Could not update user",
			"details": err.Error(),
		})
		logging.Instance.Warn("Error on updating user", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully updated user",
	})
}

// DeleteUser удаляет пользователя
// @Summary Удаление пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Produce json
// @Param userId header string true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [delete]
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    "UNAUTHORIZED",
			"message": "You're unauthorized",
		})
		return
	}

	err := h.service.DeleteUserById(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "SERVER_ERROR",
			"message": "Could not delete user",
			"details": err.Error(),
		})
		logging.Instance.Warn("Error on deleting user", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully deleted user",
	})
}

// GetAllUsers возвращает список всех пользователей
// @Summary Получение всех пользователей
// @Description Возвращает список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "SERVER_ERROR",
			"message": "Could not get users",
			"details": err.Error(),
		})
		logging.Instance.Warn("Error on getting users", err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserById получает пользователя по ID
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его ID
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUserById(ctx *gin.Context) {
	userID := ctx.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "BAD_REQUEST",
			"message": "Could not get id",
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.GetUserById(userIDInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "BAD_REQUEST",
			"message": "Could not get user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "BAD_REQUEST",
			"message": "Could not get user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
