package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"kara-bank/dto"
	"kara-bank/middlewares"
	"kara-bank/services"
	"kara-bank/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransferControllerTestSuite struct {
	suite.Suite
	ctx    context.Context
	router http.Handler
}

func TestTransferSuite(t *testing.T) {
	suite.Run(t, &TransferControllerTestSuite{})
}

func (suite *TransferControllerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	tokenMaker := utils.NewPasetoMaker("")
	validatorObj := validator.New(validator.WithRequiredStructEnabled())

	userService := services.NewUserService(testStore, tokenMaker)
	userController := NewUserController(userService, validatorObj)

	accountService := services.NewAccountService(testStore)
	accountController := NewAccountController(accountService, validatorObj)

	transferService := services.NewTransferService(testStore)
	transferController := NewTransferController(transferService, validatorObj)

	router := http.NewServeMux()

	router.HandleFunc("POST /users/register", userController.HandleRegisterUser)
	router.HandleFunc("POST /users/login", userController.HandleLoginUser)

	router.HandleFunc("POST /accounts", accountController.HandleCreateAccount)
	router.HandleFunc("GET /accounts/{id}", accountController.HandleGetAccount)
	router.HandleFunc("GET /accounts", accountController.HandleListAccounts)

	router.HandleFunc("POST /transfers", transferController.HandleCreateTransfer)

	routerWithMiddleware := middlewares.AuthMiddleware(tokenMaker, router)

	utils.SetProtectedRoutes()

	suite.router = routerWithMiddleware
}

func (suite *TransferControllerTestSuite) AfterTest(suiteName string, testName string) {
	// clear tables after every test to avoid dependencies and side effects between tests
	_, err := testStore.ClearEntriesTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearTransfersTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearAccountsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearSessionsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearUsersTable()
	require.NoError(suite.T(), err)
}

func (suite *TransferControllerTestSuite) TestCreateTransferSuccess() {
	// prepare first user and its account
	registerUserParam1 := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	accessToken1 := registerUserAndLogin(registerUserParam1, suite.router, suite.T())
	account1 := createAccount(accessToken1, "EUR", suite.router, suite.T())
	require.Equal(suite.T(), registerUserParam1.Email, account1.Owner)

	updatedAccount1, err := testStore.SetAccountBalance(suite.ctx, account1.ID, 100)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), account1.ID, updatedAccount1.ID)
	require.Equal(suite.T(), int64(100), updatedAccount1.Balance)

	// prepare second user and its account
	registerUserParam2 := &dto.RegisterUserDto{
		Email:     "Tom@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Tom",
		LastName:  "Mustermann",
	}

	accessToken2 := registerUserAndLogin(registerUserParam2, suite.router, suite.T())
	account2 := createAccount(accessToken2, "EUR", suite.router, suite.T())
	require.Equal(suite.T(), registerUserParam2.Email, account2.Owner)

	updatedAccount2, err := testStore.SetAccountBalance(suite.ctx, account2.ID, 200)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), account2.ID, updatedAccount2.ID)
	require.Equal(suite.T(), int64(200), updatedAccount2.Balance)

	// transfer money from account 1 to account 2
	transferParam := &dto.CreateTransferDto{
		FromAccountId: account1.ID,
		ToAccountId:   account2.ID,
		Amount:        100,
	}

	var body bytes.Buffer
	err = json.NewEncoder(&body).Encode(transferParam)
	require.NoError(suite.T(), err)

	request := httptest.NewRequest("POST", "/transfers", &body)
	request.AddCookie(accessToken1)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusCreated, recorder.Result().StatusCode)
}

func (suite *TransferControllerTestSuite) TestCreateTransferFailAccountAndOwnerNotMatch() {
	// prepare first user and its account
	registerUserParam1 := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	accessToken1 := registerUserAndLogin(registerUserParam1, suite.router, suite.T())
	account1 := createAccount(accessToken1, "EUR", suite.router, suite.T())
	require.Equal(suite.T(), registerUserParam1.Email, account1.Owner)

	updatedAccount1, err := testStore.SetAccountBalance(suite.ctx, account1.ID, 100)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), account1.ID, updatedAccount1.ID)
	require.Equal(suite.T(), int64(100), updatedAccount1.Balance)

	// prepare second user and its account
	registerUserParam2 := &dto.RegisterUserDto{
		Email:     "Tom@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Tom",
		LastName:  "Mustermann",
	}

	accessToken2 := registerUserAndLogin(registerUserParam2, suite.router, suite.T())
	account2 := createAccount(accessToken2, "EUR", suite.router, suite.T())
	require.Equal(suite.T(), registerUserParam2.Email, account2.Owner)

	updatedAccount2, err := testStore.SetAccountBalance(suite.ctx, account2.ID, 200)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), account2.ID, updatedAccount2.ID)
	require.Equal(suite.T(), int64(200), updatedAccount2.Balance)

	// transfer money from account 1 to account 2 but with accessToken from user 2
	transferParam := &dto.CreateTransferDto{
		FromAccountId: account1.ID,
		ToAccountId:   account2.ID,
		Amount:        100,
	}

	var body bytes.Buffer
	err = json.NewEncoder(&body).Encode(transferParam)
	require.NoError(suite.T(), err)

	request := httptest.NewRequest("POST", "/transfers", &body)
	request.AddCookie(accessToken2)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusUnauthorized, recorder.Result().StatusCode)
}
