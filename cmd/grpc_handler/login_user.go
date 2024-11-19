package gapi

import (
	"context"
	"errors"
	"kara-bank/dto"
	"kara-bank/pb"
)

func (s GrpcServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	token, respErr := s.userService.LoginUser(ctx, &dto.LoginUserDto{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: "",
		ClientIp:  "",
	})

	if respErr != nil {
		return nil, errors.New(respErr.Message)
	}

	return &pb.LoginUserResponse{
		Token: token,
	}, nil
}
