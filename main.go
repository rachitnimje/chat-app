package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Websocket server started")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the websocket server: ", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	// receive the http request from the client and upgrade to websocket protocol
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection to websocket: ", err)
	}

	fmt.Println("Client connected")

	// close the connection after exit
	//defer func(conn *websocket.Conn) {
	//	err := conn.Close()
	//	if err != nil {
	//		fmt.Println("Error closing connection:", err)
	//	}
	//}(conn)

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// read the message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
		}

		fmt.Printf("Received: %s \n", msg)

		// write the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("Error writing message: ", err)
			break
		}
	}
}
