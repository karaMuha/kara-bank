package controllers

import (
	"encoding/json"
	"kara-bank/dto"
	"kara-bank/middlewares"
	"kara-bank/services"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type TransferController struct {
	transferService services.TransferServiceInterface
	validator       *validator.Validate
}

func NewTransferController(transferService services.TransferServiceInterface, validator *validator.Validate) *TransferController {
	return &TransferController{
		transferService: transferService,
		validator:       validator,
	}
}

func (t *TransferController) HandleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateTransferDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value(middlewares.ContextUserEmailKey).(string)

	if !ok {
		http.Error(w, "Could not extract email from token", http.StatusInternalServerError)
		return
	}

	requestBody.FromUser = email
	err = t.validator.Struct(requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transfer, respErr := t.transferService.CreateTransfer(r.Context(), &requestBody)

	if respErr != nil {
		http.Error(w, respErr.Message, respErr.Status)
		return
	}

	responseJson, err := json.Marshal(&transfer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)
}
