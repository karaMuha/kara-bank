package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"kara-bank/dto"
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

type UserControllerTestSuite struct {
	suite.Suite
	ctx    context.Context
	router *http.ServeMux
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}

func (suite *UserControllerTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	userService := services.NewUserService(testStore, utils.NewPasetoMaker(""))
	userController := NewUserController(userService, validator.New(validator.WithRequiredStructEnabled()))

	router := http.NewServeMux()
	router.HandleFunc("POST /users/register", userController.HandleRegisterUser)
	router.HandleFunc("POST /users/login", userController.HandleLoginUser)

	suite.router = router
}

func (suite *UserControllerTestSuite) AfterTest(suiteName string, testName string) {
	// clear users table after every test to avoid dependencies and side effects between tests
	_, err := testStore.ClearSessionsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearUsersTable()
	require.NoError(suite.T(), err)
}

func (suite *UserControllerTestSuite) TestRegisterUserNoEmail() {
	user := &dto.RegisterUserDto{
		Email:     "",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestRegisterUserPasswordTooShort() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test123",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestRegisterUserNoFirstName() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestRegisterUserNoLastName() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 400, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestRegisterUserSuccess() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 201, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestLoginFailUserNotFound() {
	user := &dto.LoginUserDto{
		Email:    "Max@Mustermann.de",
		Password: "Test1234",
	}

	userBytes, err := json.Marshal(user)

	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)

	request := httptest.NewRequest("POST", "/users/login", body)
	request.RemoteAddr = "test"
	request.Header.Set("User-Agent", "test")
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 404, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestLoginFailWrongPassword() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 201, recorder.Result().StatusCode)

	loginRequestParam := &dto.LoginUserDto{
		Email:    "Max@Mustermann.de",
		Password: "WrongPw",
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

	require.Equal(suite.T(), 401, recorder.Result().StatusCode)
}

func (suite *UserControllerTestSuite) TestLoginSuccess() {
	user := &dto.RegisterUserDto{
		Email:     "Max@Mustermann.de",
		Password:  "Test1234",
		FirstName: "Max",
		LastName:  "Mustermann",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(userBytes)
	request := httptest.NewRequest("POST", "/users/register", body)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	require.Equal(suite.T(), 201, recorder.Result().StatusCode)

	loginRequestParam := &dto.LoginUserDto{
		Email:    "Max@Mustermann.de",
		Password: "Test1234",
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
}
