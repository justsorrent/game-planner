package handlers

import "net/http"

func HandleHealtz(w http.ResponseWriter, r *http.Request) {
	WriteJsonResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
