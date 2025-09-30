package grpc

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"userService/internal/service"

	userpb "github.com/Sayan80bayev/go-project/pkg/proto/user"
)

// UserHandler implements the gRPC UserServiceServer interface
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService *service.UserService
}

// NewUserHandler constructor
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser handles gRPC request to fetch a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	user, err := h.userService.GetUserById(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		Id:        user.ID.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),

		DeletedAt: func() *timestamppb.Timestamp {
			if user.DeletedAt != nil {
				return timestamppb.New(*user.DeletedAt)
			}
			return nil
		}(),

		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		About:     &user.About,

		DateOfBirth: func() *timestamppb.Timestamp {
			if user.DateOfBirth != nil {
				return timestamppb.New(*user.DateOfBirth)
			}
			return nil
		}(),

		AvatarUrl:       &user.AvatarURL,
		Gender:          &user.Gender,
		Location:        &user.Location,
		Socials:         user.Socials,
		NeedsCompletion: user.NeedsCompletion,
	}, nil
}
