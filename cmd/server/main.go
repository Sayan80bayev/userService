package main

import (
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "userService/docs"
	"userService/internal/bootstrap"
	"userService/internal/routes"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(r, c)

	logger.Info("Server starting on port 8080")
	err = r.Run(":" + c.Config.Port)
	if err != nil {
		logger.Errorf("Couldn't start server: %v", err)
		return
	}

}
