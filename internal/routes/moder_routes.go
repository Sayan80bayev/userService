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
		moderRoutes.PUT("role/:id", h.SetRoleById)
		moderRoutes.PUT("ban/:id", h.BanUserById)
		moderRoutes.PUT("unban/:id", h.UnBanUserById)
	}

}
