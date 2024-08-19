package controllers

import (
	"encoding/json"
	"kara-bank/dto"
	"kara-bank/services"
	"net/http"
)

type TransferController struct {
	transferService services.TransferServiceInterface
}

func NewTransferController(transferService services.TransferServiceInterface) *TransferController {
	return &TransferController{
		transferService: transferService,
	}
}

func (t *TransferController) HandleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateTransferDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)

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
