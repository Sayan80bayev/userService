package routes

import (
	"userService/internal/bootstrap"
	"userService/internal/delivery"
	"userService/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupModerRoutes(r *gin.Engine, c *bootstrap.Container) {
	ms := service.NewModerService(c.UserRepositoryImpl)
	h := delivery.NewModerHandler(ms)

	moderRoutes := r.Group("/api/v1/moder")
	{
		moderRoutes.PUT("/:id", h.SetRoleById)
		moderRoutes.PUT("/:id", h.BanUserById)
		moderRoutes.PUT("/:id", h.UnBanUserById)
	}

}
