package main

import (
	"context"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"userService/internal/bootstrap"
	"userService/internal/grpc"
	"userService/internal/routes"
)

func main() {
	logger := logging.GetLogger()

	gin.SetMode(gin.ReleaseMode)

	c, err := bootstrap.Init()
	if err != nil {
		logger.Error("Couldn't init Container", err)
	}

	// Start Gin HTTP server
	r := gin.New()
	r.Use(logging.Middleware)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupUserRoutes(r, c)

	// Start gRPC server (in its own goroutine)
	grpc.SetupGRPCServer(c)

	logger.Info("HTTP server starting on port " + c.Config.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.Consumer.Start(ctx)

	if err := r.Run(":" + c.Config.Port); err != nil {
		logger.Errorf("Couldn't start Gin server: %v", err)
		return
	}

	defer c.Consumer.Close()
}
