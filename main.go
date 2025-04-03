package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var clients = make(map[*websocket.Conn]struct{})
var broadcastChannel = make(chan []byte)
var mutex = &sync.Mutex{}

func main() {
	// each ws connection runs in its own goroutine
	http.HandleFunc("/ws", wsHandler)
	// we create a separate goroutine for broadcasting messages
	go handleMessages()

	fmt.Println("Websocket server started on port 8080")

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

	// mutex lock to prevent race condition
	mutex.Lock()
	clients[conn] = struct{}{}
	mutex.Unlock()

	// continuously check for new message from the client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		// add the new message to the broadcastChannel to be broadcast to every connection
		fmt.Printf("Received msg: %s\n", msg)
		broadcastChannel <- msg
	}
}

func handleMessages() {
	// continuously check for new message to the broadcastChannel
	for {
		message := <-broadcastChannel
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
