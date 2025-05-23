package main

import (
	"github.com/joho/godotenv"
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/config"
	"github.com/rachitnimje/chat-app/internal/database"
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/server"
	"github.com/rachitnimje/chat-app/pkg/routes"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	log.Printf(".env file loaded")

	auth.InitJWTSecret()
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
	port := os.Getenv("PORT")
	log.Printf("websocket server started on port %s\n", port)

	// start the web server
	err = server.StartHTTPServer(router, port)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
