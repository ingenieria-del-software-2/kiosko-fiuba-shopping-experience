package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of API error responses
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteErrorResponse writes an error response to the response writer
func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	response := ErrorResponse{
		Status:  status,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
