package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func WriteSuccessResponse(writer http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	WriteJSONResponse(writer, response)
}

func WriteErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	response := Response{
		StatusCode: statusCode,
		Message:    message,
	}
	WriteJSONResponse(writer, response)
}

func WriteJSONResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
