package main

import (
	"github.com/gin-gonic/gin"
	"userService/internal/bootstrap"
	"userService/pkg/logging"
)

var logger = logging.GetLogger()

func main() {
	gin.SetMode(gin.ReleaseMode)
	bs, err := bootstrap.Init()
	if err != nil {
		logger.Error("Couldn't init Container", err)
	}

	r := gin.New()
	r.Use(logging.Middleware)

	logger.Info("Server starting on port 8080")
	r.Run(":" + bs.Config.Port)

}
