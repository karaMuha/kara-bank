package rest

import (
	"encoding/json"
	"kara-bank/dto"
	"kara-bank/middlewares"
	"kara-bank/services"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type AccountController struct {
	accountService services.AccountServiceInterface
	validator      *validator.Validate
}

func NewAccountController(accountService services.AccountServiceInterface, validator *validator.Validate) *AccountController {
	return &AccountController{
		accountService: accountService,
		validator:      validator,
	}
}

func (a *AccountController) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateAccountDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value(middlewares.ContextUserEmailKey).(string)

	if !ok {
		http.Error(w, "Could not convert email from token to string", http.StatusInternalServerError)
		return
	}

	requestBody.Owner = email
	err = a.validator.Struct(requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, respErr := a.accountService.CreateAccount(r.Context(), &requestBody)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	responseJson, err := json.Marshal(&account)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}

func (a *AccountController) HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value(middlewares.ContextUserEmailKey).(string)

	if !ok {
		http.Error(w, "Could not convert email from token to string", http.StatusInternalServerError)
		return
	}

	role, ok := r.Context().Value(middlewares.ContextUserRoleKey).(string)

	if !ok {
		http.Error(w, "Could not convert email from token to string", http.StatusInternalServerError)
		return
	}

	account, respErr := a.accountService.GetAccount(r.Context(), int64(id), email, role)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	responseJson, err := json.Marshal(&account)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (a *AccountController) HandleListAccounts(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.ListAccountsDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.validator.Struct(requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value(middlewares.ContextUserRoleKey).(string)

	if !ok {
		http.Error(w, "Could not convert email from token to string", http.StatusInternalServerError)
		return
	}

	accounts, respErr := a.accountService.ListAccounts(r.Context(), &requestBody, role)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	responseJson, err := json.Marshal(&accounts)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
