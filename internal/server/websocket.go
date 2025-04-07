package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var clients = make(map[*websocket.Conn]struct{})
var broadcastChannel = make(chan []byte)
var mutex = &sync.Mutex{}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func WSHandler(writer http.ResponseWriter, request *http.Request) {
	// receive the http request from the client and upgrade to websocket protocol
	conn, err := upgrade.Upgrade(writer, request, nil)
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

func StartWSServer() {
	// we create a separate goroutine for broadcasting messages
	go handleMessages()
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
