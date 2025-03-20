package routes

import (
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/internal/delivery"
	"userService/internal/service"
)

func SetupUserRoutes(r *gin.Engine, c *bootstrap.Container) {
	s := service.NewUserService(c.UserRepositoryImpl, c.FileService, c.Producer)
	h := delivery.NewUserHandler(s)

	routes := r.Group("api/v1/users")
	{
		routes.GET("/", h.GetAllUsers)
		routes.GET("/:id", h.GetUserById)
		routes.GET("/:username", h.GetUserByUsername)
		routes.PUT("/:id", h.UpdateUser)
		routes.DELETE("/:id", h.DeleteUser)
	}
}
