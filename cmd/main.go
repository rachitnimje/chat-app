package main

import (
	"github.com/rachitnimje/chat-app/internal/config"
	"github.com/rachitnimje/chat-app/internal/database"
	"github.com/rachitnimje/chat-app/internal/server"
	"github.com/rachitnimje/chat-app/pkg/routes"
	"gorm.io/gorm"
	"log"
)

func main() {
	cfg := config.DefaultConfig()

	// connect to database
	var DB *gorm.DB
	dbConfig := cfg.Database
	err := database.ConnectDB(DB, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	// create a router along with routes
	router := routes.NewRouter()

	// start the websocket server
	server.StartWSServer()
	port := 8080
	log.Printf("websocket server started on port %d\n", port)

	// start the web server at port 8080
	err = server.StartHTTPServer(router, port)
	if err != nil {
		log.Fatal("error starting http server: ", err)
	}
}
