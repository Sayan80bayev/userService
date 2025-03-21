package routes

import (
	"userService/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, c *bootstrap.Container) {
	SetupUserRoutes(r, c)
	SetupModerRoutes(r, c)
}
