package server

import (
	"kara-bank/middlewares"
	rest "kara-bank/rest_handler"
	"kara-bank/services"
	"kara-bank/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func InitHttpServer(
	port string,
	userService services.UserServiceInterface,
	accountService services.AccountServiceInterface,
	transferService services.TransferServiceInterface,
	tokenMaker utils.TokenMaker,
) *http.Server {
	// init validator
	validator := validator.New(validator.WithRequiredStructEnabled())

	// init controller layer
	userController := rest.NewUserController(userService, validator)
	accountController := rest.NewAccountController(accountService, validator)
	transferController := rest.NewTransferController(transferService, validator)

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /users/register", userController.HandleRegisterUser)
	router.HandleFunc("POST /users/login", userController.HandleLoginUser)

	router.HandleFunc("POST /accounts", accountController.HandleCreateAccount)
	router.HandleFunc("GET /accounts/{id}", accountController.HandleGetAccount)
	router.HandleFunc("GET /accounts", accountController.HandleListAccounts)

	router.HandleFunc("POST /transfers", transferController.HandleCreateTransfer)

	// init protected routes
	utils.SetProtectedRoutes()

	return &http.Server{
		Addr:    port,
		Handler: middlewares.AuthMiddleware(tokenMaker, router),
	}
}
