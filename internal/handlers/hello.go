package handlers

import (
	"github.com/rachitnimje/chat-app/internal/utils"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Name":     "Rachit",
		"username": "rachitnimje",
	}
	utils.WriteJSONResponse(w, http.StatusOK, data)
}
