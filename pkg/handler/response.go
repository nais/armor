package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error-message"`
}

func HttpError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	er := errorResponse{
		Status:       status,
		ErrorMessage: message,
	}
	err := json.NewEncoder(w).Encode(er)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
}

func response(w http.ResponseWriter, response interface{}) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
}
