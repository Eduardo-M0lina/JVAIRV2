package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse representa una respuesta de error en JSON
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse representa una respuesta exitosa en JSON
type SuccessResponse struct {
	Message string `json:"message"`
}

// DataResponse representa una respuesta con datos en JSON
type DataResponse struct {
	Data interface{} `json:"data"`
}

// RespondWithError envía una respuesta de error en formato JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// RespondWithSuccess envía una respuesta exitosa en formato JSON
func RespondWithSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SuccessResponse{Message: message})
}

// RespondWithJSON envía una respuesta con datos en formato JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
