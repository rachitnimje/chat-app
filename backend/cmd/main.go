package main

import (
	"github.com/joho/godotenv"
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/config"
	"github.com/rachitnimje/chat-app/internal/database"
	"github.com/rachitnimje/chat-app/internal/handlers"
	"github.com/rachitnimje/chat-app/internal/models"
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

	// migrate the database
	if err := models.Migrate(db); err != nil {
		log.Fatal("error migrating database: ", err)
	}

	// initialize websocket with database
	server.InitWebsocket(db)

	// handlers
	authHandler := handlers.NewAuthHandler(db)
	roomHandler := handlers.NewRoomHandler(db)

	// routes
	router := routes.NewRouter(authHandler, roomHandler)

	// start the websocket server
	server.StartWSServer()
	port := os.Getenv("PORT")

	log.Printf("websocket server starting on port %s\n", port)

	// start the web server
	err = server.StartHTTPServer(router, port)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
