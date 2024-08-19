package controllers

import (
	"encoding/json"
	"kara-bank/dto"
	"kara-bank/services"
	"net/http"
	"time"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.RegisterUserDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, respErr := u.userService.RegisterUser(r.Context(), &requestBody)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	responseJson, err := json.Marshal(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}

func (u *UserController) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.LoginUserDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, respErr := u.userService.LoginUser(r.Context(), &requestBody)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "paseto",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * time.Minute),
	})

	w.WriteHeader(http.StatusOK)
}
