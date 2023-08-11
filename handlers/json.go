package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	if status >= 500 {
		log.Printf("Error: %v\n", message)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, status, errResponse{Error: message})
}
