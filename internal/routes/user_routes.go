package routes

import (
	"github.com/Sayan80bayev/go-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/internal/delivery"
	"userService/internal/service"
)

func SetupUserRoutes(r *gin.Engine, c *bootstrap.Container) {
	repo := c.UserRepository
	s := service.NewUserService(repo, c.FileStorage, c.Producer, c.Redis)
	h := delivery.NewUserHandler(s)

	routes := r.Group("api/v1/users")
	{
		routes.GET("", h.GetAllUsers)
		routes.GET("/:id", h.GetUserById)
		// routes.GET("/", h.GetUserByUsername)
	}

	authRoutes := r.Group("api/v1/users", middleware.AuthMiddleware(c.JWKSUrl))
	{
		authRoutes.PUT("/:id", h.UpdateUser)
		authRoutes.DELETE("/:id", h.DeleteUser)
	}
}
