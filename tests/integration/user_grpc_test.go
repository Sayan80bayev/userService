package integration

import (
	"context"
	"testing"
	"time"

	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/Sayan80bayev/go-project/pkg/proto/user"
	"userService/internal/events"
)

func newGRPCClient(t *testing.T) userpb.UserServiceClient {
	conn, err := grpc.NewClient(
		"localhost:50052", // adjust if needed
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	logging.GetLogger().Infof("Created gRPC Client")
	require.NoError(t, err)
	return userpb.NewUserServiceClient(conn)
}

func TestGetUser_GRPC(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	logger := logging.GetLogger()
	logger.Infof("testing gRPC GetUser %s", userID)

	// --- Step 1: Produce Kafka user.created event ---
	payload := events.UserCreatedPayload{
		UserID:    userID,
		Firstname: "Sayan",
		Lastname:  "Seksenbayev",
		Email:     "sayan123serv@gmail.com",
	}

	err := container.Producer.Produce(ctx, events.UserCreated, payload)
	require.NoError(t, err)

	// --- Step 2: Wait for event to be consumed and user stored ---
	time.Sleep(5 * time.Second)

	client := newGRPCClient(t)

	// --- Step 3: Fetch created user via gRPC ---
	res, err := client.GetUser(ctx, &userpb.GetUserRequest{UserId: userID.String()})
	require.NoError(t, err, "GetUser should return user without error")
	require.Equal(t, "Sayan", res.Firstname)
	require.Equal(t, "Seksenbayev", res.Lastname)
	require.Equal(t, "sayan123serv@gmail.com", res.Email)
}
