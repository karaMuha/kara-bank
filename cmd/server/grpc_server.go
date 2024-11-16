package server

import (
	"kara-bank/pb"
	"kara-bank/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GrpcServer struct {
	pb.UnimplementedKaraBankServer
	db         *pgxpool.Pool
	tokenMaker utils.TokenMaker
}

func InitGrpcHandler(connPool *pgxpool.Pool, tokenMaker utils.TokenMaker) *GrpcServer {
	// init repository layer
	// store := db.NewStore(connPool)

	// init service layer
	/* userService := services.NewUserService(store, tokenMaker)
	accountService := services.NewAccountService(store)
	transferService := services.NewTransferService(store) */

	return &GrpcServer{}
}
