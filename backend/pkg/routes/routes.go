package routes

import (
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/server"
	"net/http"
)

func NewRouter(authHandler *handlers.AuthHandler, roomHandler *handlers.RoomHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/ws", server.WSHandler)

	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/register", authHandler.Register)

	// PROTECTED ROUTES
	router.Handle("POST /rooms", auth.Middleware(http.HandlerFunc(roomHandler.CreateRoom)))
	router.Handle("GET /rooms", auth.Middleware(http.HandlerFunc(roomHandler.GetRooms)))
	router.Handle("GET /messages", auth.Middleware(http.HandlerFunc(roomHandler.GetMessages)))

	return router
}
