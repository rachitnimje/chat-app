package handlers

import (
	"encoding/json"
	"github.com/rachitnimje/chat-app/internal/models"
	"github.com/rachitnimje/chat-app/internal/utils"
	"gorm.io/gorm"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthHandler acts as container for dependencies, injects the dependencies through the struct
type AuthHandler struct {
	DB *gorm.DB
}

// NewAuthHandler is a constructor function which initializes handler with all required dependencies
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Login is a method attached to AuthHandler struct, variable h gives you access to the fields of struct
// (h *AuthHandler) is a pointer method receiver, this way you can modify the receiver instance and access its fields
// (h AuthHandler) is a value method receiver, this way you receive the copy of the original
func (h AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteErrorResponse(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// validation
	var loginRequest LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "Username and password are required")
		return
	}

	var user models.User
	result := h.DB.Where("username = ? AND password = ?", loginRequest.Username, loginRequest.Password).First(&user)
	if result.Error != nil {
		utils.WriteErrorResponse(writer, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response := LoginResponse{
		Token: "login successfully",
	}

	utils.WriteSuccessResponse(writer, http.StatusOK, "Logged in successfully", response)
}
