package main

import (
	"github.com/rachitnimje/chat-app/internal/config"
	"github.com/rachitnimje/chat-app/internal/database"
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/server"
	"github.com/rachitnimje/chat-app/pkg/routes"
	"log"
)

func main() {
	cfg := config.DefaultConfig()

	// connect to database
	dbConfig := cfg.Database
	db, err := database.ConnectDB(
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
	)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
		return
	}

	// handlers
	authHandler := handlers.NewAuthHandler(db)

	// routes
	router := routes.NewRouter(authHandler)

	// start the websocket server
	server.StartWSServer()
	port := 8080
	log.Printf("websocket server started on port %d\n", port)

	// start the web server
	err = server.StartHTTPServer(router, port)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
