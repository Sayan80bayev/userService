package routes

import (
	"github.com/Sayan80bayev/go-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/internal/delivery"
)

func SetupUserRoutes(r *gin.Engine, c *bootstrap.Container) {
	h := delivery.NewUserHandler(c.UserService)

	routes := r.Group("api/v1/users")
	{
		routes.GET("", h.GetAllUsers)
		routes.GET("/:id", h.GetUserById)
		// routes.GET("/", h.GetUserByUsername)
	}

	authRoutes := r.Group("api/v1/users", middleware.AuthMiddleware(c.JWKSUrl))
	{
		authRoutes.DELETE("/:id", h.DeleteUser)
		authRoutes.PUT("/:id", h.UpdateUser)
	}
}
