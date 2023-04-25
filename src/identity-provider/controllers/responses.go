package controllers

import (
	"encoding/json"
	"net/http"
)

type ErrorDescription struct {
	Field string `json:"filed"`
	Error string `json:"error"`
}
type validationErrorResponse struct {
	Message string             `json:"message"`
	Errors  []ErrorDescription `json:"errors"`
}

func ValidationErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)

	resp := &validationErrorResponse{message, []ErrorDescription{}}
	json.NewEncoder(w).Encode(resp)
}

func BadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

func JsonSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Response-Code", "00")
	w.Header().Set("Response-Desc", "Success")

	json.NewEncoder(w).Encode(data)
}