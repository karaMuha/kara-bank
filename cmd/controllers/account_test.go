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
	// clear users table after every test to avoid dependencies and side effects between tests
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
	accessTokenCookie := suite.registerUserAndLogin(&dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	})

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
	accessTokenCookie := suite.registerUserAndLogin(&dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	})

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
	accessTokenCookie := suite.registerUserAndLogin(&registerUserParam)

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

	// get created account
	var createdAccount db.Account
	err = json.NewDecoder(recorder.Result().Body).Decode(&createdAccount)

	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), createdAccount.ID)
	require.Equal(suite.T(), registerUserParam.Email, createdAccount.Owner)
}

func (suite *AccountControllerTestSuite) registerUserAndLogin(arg *dto.RegisterUserDto) *http.Cookie {
	// register user
	userBytes, err := json.Marshal(arg)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 201, recorder.Result().StatusCode)

	// login with registered user
	loginRequestParam := &dto.LoginUserDto{
		Email:    arg.Email,
		Password: arg.Password,
	}

	loginRequestParamBytes, err := json.Marshal(loginRequestParam)
	if err != nil {
		log.Fatal(err)
	}

	requestBody := bytes.NewReader(loginRequestParamBytes)
	request = httptest.NewRequest("POST", "/users/login", requestBody)
	request.RemoteAddr = "test"
	request.Header.Set("User-Agent", "test")
	recorder = httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 200, recorder.Result().StatusCode)

	accessTokenCookie := getCookie(recorder.Result().Cookies())
	require.NotNil(suite.T(), accessTokenCookie)

	return accessTokenCookie
}

func getCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie
		}
	}

	return nil
}
