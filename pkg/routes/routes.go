package routes

import (
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/server"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /hello", handlers.HelloHandler)
	router.HandleFunc("/ws", server.WSHandler)

	return router
}
