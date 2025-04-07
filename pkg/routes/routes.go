package routes

import (
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/server"
	"net/http"
)

func NewRouter(authHandler *handlers.AuthHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /hello", handlers.HelloHandler)
	router.HandleFunc("/ws", server.WSHandler)
	router.HandleFunc("POST /auth/login", authHandler.Login)

	return router
}
