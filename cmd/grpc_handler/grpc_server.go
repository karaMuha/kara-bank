package gapi

import (
	"kara-bank/pb"
	"kara-bank/services"
)

type GrpcServer struct {
	pb.UnimplementedKaraBankServer
	userService    services.UserServiceInterface
	accountService services.AccountServiceInterface
	transerService services.TransferServiceInterface
}

func InitGrpcHandler(
	userService services.UserServiceInterface,
	accountService services.AccountServiceInterface,
	transferService services.TransferServiceInterface,
) *GrpcServer {
	return &GrpcServer{
		userService:    userService,
		accountService: accountService,
		transerService: transferService,
	}
}
