package grpc

import (
	"net"
	"userService/internal/bootstrap"
	userpb "userService/proto"

	"github.com/Sayan80bayev/go-project/pkg/logging"
	"google.golang.org/grpc"
)

// SetupGRPCServer starts the gRPC server in a goroutine
func SetupGRPCServer(c *bootstrap.Container) {
	logger := logging.GetLogger()

	s := c.UserService
	h := NewUserHandler(s)

	go func() {
		lis, err := net.Listen("tcp", ":"+c.Config.GrpcPort)
		if err != nil {
			logger.Fatalf("failed to listen on :%s %v", c.Config.GrpcPort, err)
		}

		grpcServer := grpc.NewServer()
		userpb.RegisterUserServiceServer(grpcServer, h)

		logger.Infof("gRPC server started on %s", c.Config.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
}
