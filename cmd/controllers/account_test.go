package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	db "kara-bank/db/repositories"
	"kara-bank/dto"
	"kara-bank/middlewares"
	"kara-bank/services"
	"kara-bank/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountControllerTestSuite struct {
	suite.Suite
	ctx    context.Context
	router http.Handler
}

func TestAccountControllerTestSuite(t *testing.T) {
	suite.Run(t, &AccountControllerTestSuite{})
}

func (suite *AccountControllerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	tokenMaker := utils.NewPasetoMaker("")

	userService := services.NewUserService(testStore, tokenMaker)
	userController := NewUserController(userService, validator.New(validator.WithRequiredStructEnabled()))

	accountService := services.NewAccountService(testStore)
	accountController := NewAccountController(accountService, validator.New(validator.WithRequiredStructEnabled()))

	router := http.NewServeMux()

	router.HandleFunc("POST /users/register", userController.HandleRegisterUser)
	router.HandleFunc("POST /users/login", userController.HandleLoginUser)

	router.HandleFunc("POST /accounts", accountController.HandleCreateAccount)
	router.HandleFunc("GET /accounts/{id}", accountController.HandleGetAccount)
	router.HandleFunc("GET /accounts", accountController.HandleListAccounts)

	routerWithMiddleware := middlewares.AuthMiddleware(tokenMaker, router)

	utils.SetProtectedRoutes()

	suite.router = routerWithMiddleware
}

func (suite *AccountControllerTestSuite) AfterTest(suiteName string, testName string) {
	// clear tables after every test to avoid dependencies and side effects between tests
	_, err := testStore.ClearAccountsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearSessionsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearUsersTable()
	require.NoError(suite.T(), err)
}

func (suite *AccountControllerTestSuite) TestCreateAccountNotLoggedIn() {
	createAccountParam := &dto.CreateAccountDto{
		Currency: "EUR",
	}

	createAccountParamBytes, err := json.Marshal(createAccountParam)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(createAccountParamBytes)
	request := httptest.NewRequest("POST", "/accounts", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusUnauthorized, recorder.Result().StatusCode)
}

func (suite *AccountControllerTestSuite) TestCreateAccountSuccess() {
	accessTokenCookie := registerUserAndLogin(&dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}, suite.router, suite.T())
	require.NotNil(suite.T(), accessTokenCookie)

	// create account with logged in user
	createAccountParam := &dto.CreateAccountDto{
		Currency: "EUR",
	}

	createAccountParamBytes, err := json.Marshal(createAccountParam)
	if err != nil {
		log.Fatal(err)
	}

	requestBody := bytes.NewReader(createAccountParamBytes)
	request := httptest.NewRequest("POST", "/accounts", requestBody)
	request.AddCookie(accessTokenCookie)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusCreated, recorder.Result().StatusCode)
}

func (suite *AccountControllerTestSuite) TestGetOneAccountNotFound() {
	accessTokenCookie := registerUserAndLogin(&dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}, suite.router, suite.T())
	require.NotNil(suite.T(), accessTokenCookie)

	request := httptest.NewRequest("GET", "/accounts/123", nil)
	request.AddCookie(accessTokenCookie)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusNotFound, recorder.Result().StatusCode)
}

func (suite *AccountControllerTestSuite) TestGetOneAccountSuccess() {
	registerUserParam := dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}
	accessTokenCookie := registerUserAndLogin(&registerUserParam, suite.router, suite.T())
	require.NotNil(suite.T(), accessTokenCookie)

	// create account with logged in user
	createAccountParam := &dto.CreateAccountDto{
		Currency: "EUR",
	}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(createAccountParam)
	require.NoError(suite.T(), err)

	request := httptest.NewRequest("POST", "/accounts", &body)
	request.AddCookie(accessTokenCookie)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusCreated, recorder.Result().StatusCode)

	// get created account
	var createdAccount db.Account
	err = json.NewDecoder(recorder.Result().Body).Decode(&createdAccount)

	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), createdAccount.ID)
	require.Equal(suite.T(), registerUserParam.Email, createdAccount.Owner)

	endpoint := "/accounts/" + strconv.Itoa(int(createdAccount.ID))

	request = httptest.NewRequest("GET", endpoint, &body)
	request.AddCookie(accessTokenCookie)
	recorder = httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), http.StatusOK, recorder.Result().StatusCode)
}

// helper function for test suits that need accounts
func createAccount(accessToken *http.Cookie, currency string, router http.Handler, t *testing.T) *db.Account {
	createAccountParam := &dto.CreateAccountDto{
		Currency: currency,
	}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(createAccountParam)
	require.NoError(t, err)

	request := httptest.NewRequest("POST", "/accounts", &body)
	request.AddCookie(accessToken)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusCreated, recorder.Result().StatusCode)

	var createdAccount db.Account
	err = json.NewDecoder(recorder.Result().Body).Decode(&createdAccount)
	require.NoError(t, err)

	return &createdAccount
}
