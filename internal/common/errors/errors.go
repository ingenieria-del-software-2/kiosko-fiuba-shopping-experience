package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteErrorResponse writes an error response with the given status code and message
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	// Map status codes to error codes
	var code string
	switch statusCode {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusNotFound:
		code = "NOT_FOUND"
	case http.StatusForbidden:
		code = "FORBIDDEN"
	case http.StatusConflict:
		code = "CONFLICT"
	case http.StatusUnprocessableEntity:
		code = "UNPROCESSABLE_ENTITY"
	default:
		code = "INTERNAL_SERVER_ERROR"
	}

	// Create the error response
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
