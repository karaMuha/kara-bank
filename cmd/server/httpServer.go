package server

import (
	"kara-bank/controllers"
	db "kara-bank/db/repositories"
	"kara-bank/middlewares"
	"kara-bank/services"
	"kara-bank/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitHttpServer(port string, connPool *pgxpool.Pool, tokenMaker utils.TokenMaker) *http.Server {
	// init validator
	validator := validator.New(validator.WithRequiredStructEnabled())
	// init repository layer
	store := db.NewStore(connPool)

	// init service layer
	userService := services.NewUserService(store, tokenMaker)
	accountService := services.NewAccountService(store)
	transferService := services.NewTransferService(store)

	// init controller layer
	userController := controllers.NewUserController(userService, validator)
	accountController := controllers.NewAccountController(accountService, validator)
	transferController := controllers.NewTransferController(transferService, validator)

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /users/register", userController.HandleRegisterUser)
	router.HandleFunc("POST /users/login", userController.HandleLoginUser)

	router.HandleFunc("POST /accounts", accountController.HandleCreateAccount)
	router.HandleFunc("GET /accounts/{id}", accountController.HandleGetAccount)
	router.HandleFunc("GET /accounts", accountController.HandleListAccounts)

	router.HandleFunc("POST /transfers", transferController.HandleCreateTransfer)

	return &http.Server{
		Addr:    port,
		Handler: middlewares.AuthMiddleware(tokenMaker, router),
	}
}
