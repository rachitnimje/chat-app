package routes

import (
	"github.com/rachitnimje/chat-app/internal"
	"github.com/rachitnimje/chat-app/internal/utils"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /hello", helloHandler)
	router.HandleFunc("/ws", internal.WSHandler)

	return router
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Name":     "Rachit",
		"username": "rachitnimje",
	}
	utils.WriteJSONResponse(w, http.StatusOK, data)
}
