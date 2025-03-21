package main

import (
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/internal/routes"
	"userService/pkg/logging"
)

var logger = logging.GetLogger()

func main() {
	gin.SetMode(gin.ReleaseMode)
	c, err := bootstrap.Init()
	if err != nil {
		logger.Error("Couldn't init Container", err)
	}

	r := gin.New()
	r.Use(logging.Middleware)
	routes.SetupRoutes(r, c)

	logger.Info("Server starting on port 8080")
	err = r.Run(":" + c.Config.Port)
	if err != nil {
		logger.Errorf("Couldn't start server: %v", err)
		return
	}

}
