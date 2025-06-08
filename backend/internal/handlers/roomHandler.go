package handlers

import (
	"encoding/json"
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/models"
	"github.com/rachitnimje/chat-app/internal/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RoomHandler struct {
	db *gorm.DB
}

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

func (h *RoomHandler) CreateRoom(writer http.ResponseWriter, request *http.Request) {
	// retrieve the username from context
	username, err := auth.GetUsernameFromContext(request.Context())
	if err != nil {
		utils.WriteErrorResponse(writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// retrieve the user from db using username
	var user models.User
	if err = h.db.Where("username = ?", username).First(&user).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusNotFound, "User not found")
		return
	}

	// decode the create room request from request object
	var roomReq CreateRoomRequest
	if err := json.NewDecoder(request.Body).Decode(&roomReq).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// create a new room
	room := &models.Room{
		Name:        roomReq.Name,
		Description: roomReq.Description,
		CreatedBy:   user.ID,
	}

	if err = h.db.Create(&room).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusInternalServerError, "Failed to create the room")
		return
	}
	utils.WriteSuccessResponse(writer, http.StatusCreated, "Room created successfully", room)
}

func (h *RoomHandler) GetRooms(writer http.ResponseWriter, request *http.Request) {
	var rooms []models.Room
	if err := h.db.Find(&rooms).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusNotFound, "Failed to fetch rooms")
		return
	}

	utils.WriteSuccessResponse(writer, http.StatusOK, "Rooms fetched successfully", rooms)
}

func (h *RoomHandler) GetMessages(writer http.ResponseWriter, request *http.Request) {
	// get the room id from the url
	roomIDStr := request.URL.Query().Get("room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "Invalid room ID")
		return
	}

	// query the db for messages with specified room id
	var messages []models.Message
	if err = h.db.Where("room_id = ?", roomID).Preload("User").Find(&messages).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusInternalServerError, "Failed to fetch messages")
		return
	}

	utils.WriteSuccessResponse(writer, http.StatusOK, "Messages fetched successfully", messages)
}
