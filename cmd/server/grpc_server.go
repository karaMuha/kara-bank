package server

import "google.golang.org/grpc"

func InitGrpcServer() *grpc.Server {
	server := grpc.NewServer()

	return server
}
