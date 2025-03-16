package main

import (
	"github.com/gin-gonic/gin"
	"userService/pkg/logging"
)

var logger = logging.GetLogger()

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(logging.Middleware)

	logger.Info("Server starting on port 8080")
	r.Run(":" + "8080")

}
