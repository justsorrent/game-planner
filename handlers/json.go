package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(res)
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
