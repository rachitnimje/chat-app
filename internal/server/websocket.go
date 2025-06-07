package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/models"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

// Client stores the data of the connected client/user
type Client struct {
	conn     *websocket.Conn
	userID   uint
	roomID   uint
	username string
}

// MessagePayload is the payload that we will receive from the websocket connection
type MessagePayload struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	RoomID  uint   `json:"room_id"`
	Token   string `json:"token,omitempty"`
}

// MessageData is the payload that will be broadcast to the channel
type MessageData struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomID   uint   `json:"room_id"`
	Client   *Client
}

var clients = make(map[*Client]struct{})
var broadcastChannel = make(chan MessageData)
var rooms = make(map[uint]map[*Client]struct{})
var mutex = &sync.Mutex{}
var db *gorm.DB

var upgrade = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func InitWebsocket(database *gorm.DB) {
	db = database
	if rooms == nil {
		rooms = make(map[uint]map[*Client]struct{})
	}
}

func WSHandler(writer http.ResponseWriter, request *http.Request) {
	// receive the http request from the client and upgrade to websocket protocol
	conn, err := upgrade.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection to websocket: ", err)
		return
	}

	client := &Client{conn: conn}

	// mutex lock to prevent race condition
	mutex.Lock()
	clients[client] = struct{}{}
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, client)

		mutex.Unlock()
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err)
			return
		}
	}()

	// continuously check for new message from the client
	for {
		var msg MessagePayload
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}

		switch msg.Type {
		case "join_room":
			handleJoinRoom(client, msg)
		case "message":
			handleMessage(client, msg)
		}
	}
}

// handleMessage handles the incoming message to the websocket connection
func handleMessage(client *Client, msg MessagePayload) {
	// create a struct which will be stored in the database using the models.Message struct
	message := models.Message{
		Content: msg.Content,
		RoomID:  msg.RoomID,
		UserID:  client.userID,
	}

	// save the message struct to the database
	if err := db.Create(&message).Error; err != nil {
		fmt.Println("Error saving message: ", err)
		return
	}

	// broadcast the message to broadcastChannel
	broadcastChannel <- MessageData{
		Type:     msg.Type,
		Content:  msg.Content,
		Username: client.username,
		RoomID:   msg.RoomID,
		Client:   client,
	}
}

// handleJoinRoom
func handleJoinRoom(client *Client, msg MessagePayload) {
	// verify the token of the client
	token, err := auth.VerifyToken(msg.Token)
	if err != nil {
		if err := client.conn.WriteJSON(map[string]string{"error": "Invalid token"}); err != nil {
			fmt.Println("Error writing to connection: ", err)
			return
		}
		return
	}

	// extract username from the token
	username, err := auth.ExtractUsername(token)
	if err != nil {
		if err := client.conn.WriteJSON(map[string]string{"error": "Invalid token"}); err != nil {
			fmt.Println("Error writing to connection: ", err)
			return
		}
		return
	}

	// fetch the user from the database using username
	var user models.User
	if err = db.Where("username = ?", username).First(&user).Error; err != nil {
		if err := client.conn.WriteJSON(map[string]string{"error": "User not found"}); err != nil {
			fmt.Println("Error writing to connection: ", err)
			return
		}
		return
	}

	client.userID = user.ID
	client.username = username
	client.roomID = msg.RoomID

	mutex.Lock()
	if rooms[msg.RoomID] == nil {
		rooms[msg.RoomID] = make(map[*Client]struct{})
	}
	rooms[msg.RoomID][client] = struct{}{}
	mutex.Unlock()

	// Send confirmation
	client.conn.WriteJSON(map[string]interface{}{
		"type":     "joined_room",
		"room_id":  msg.RoomID,
		"username": username,
	})
}

// StartWSServer creates a separate goroutine for every websocket connection
func StartWSServer() {
	go handleMessages()
}

// handleMessages continuously checks for new message to the broadcastChannel
func handleMessages() {
	for {
		msgData := <-broadcastChannel
		mutex.Lock()

		if roomClients, exists := rooms[msgData.RoomID]; exists {
			for client := range roomClients {
				err := client.conn.WriteJSON(map[string]interface{}{
					"type":     msgData.Type,
					"content":  msgData.Content,
					"username": msgData.Username,
					"room_id":  msgData.RoomID,
				})
				if err != nil {
					err := client.conn.Close()
					if err != nil {
						fmt.Println("Error closing the connection: ", err)
						return
					}
					delete(roomClients, client)
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
}
