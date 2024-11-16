package gapi

import (
	"context"
	"errors"
	"kara-bank/dto"
	"kara-bank/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s GrpcServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user, respErr := s.userService.RegisterUser(ctx, &dto.RegisterUserDto{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if respErr != nil {
		return nil, errors.New(respErr.Message)
	}

	return &pb.RegisterUserResponse{
		User: &pb.User{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UserRole:  user.UserRole,
		},
	}, nil
}
