package routes

import (
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/internal/delivery"
	"userService/internal/pkg/middleware"
	"userService/internal/service"
)

func SetupUserRoutes(r *gin.Engine, c *bootstrap.Container) {
	s := service.NewUserService(c.UserRepositoryImpl, c.FileService, c.Producer)
	h := delivery.NewUserHandler(s)

	routes := r.Group("api/v1/users")
	{
		routes.GET("", h.GetAllUsers)
		routes.GET("/:id", h.GetUserById)
		// routes.GET("/", h.GetUserByUsername)
	}

	authRoutes := r.Group("api/v1/users", middleware.AuthMiddleware(c.Config.JWTSecret))
	{
		authRoutes.PUT("/:id", h.UpdateUser)
		authRoutes.DELETE("/:id", h.DeleteUser)
		authRoutes.PATCH("/pwd/:id", h.ChangePassword)
	}
}
