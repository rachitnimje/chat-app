package handlers

import (
	"encoding/json"
	"github.com/rachitnimje/chat-app/internal/auth"
	"github.com/rachitnimje/chat-app/internal/models"
	"github.com/rachitnimje/chat-app/internal/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AuthHandler acts as container for dependencies, injects the dependencies through the struct
type AuthHandler struct {
	DB *gorm.DB
}

// NewAuthHandler is a constructor function which initializes handler with required dependencies
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

	// receive the login request credentials from the request
	var loginRequest LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validation
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

	token, err := auth.CreateToken(loginRequest.Username)
	if err != nil {
		log.Fatal("Error creating token: ", err)
		return
	}

	response := LoginResponse{
		Token: token,
	}

	utils.WriteSuccessResponse(writer, http.StatusOK, "Logged in successfully", response)
}

func (h AuthHandler) Register(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteErrorResponse(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// decode the json from the request
	var registerRequest RegisterRequest
	if err := json.NewDecoder(request.Body).Decode(&registerRequest); err != nil {
		utils.WriteErrorResponse(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// validate the request data
	if registerRequest.Username != "" || registerRequest.Name != "" || registerRequest.Password != "" {
		utils.WriteErrorResponse(writer, http.StatusBadRequest, "All fields are required")
		return
	}

	var existingUser models.User
	result := h.DB.Where("username=?", registerRequest.Username).First(&existingUser)
	if result.Error != nil {
		utils.WriteErrorResponse(writer, http.StatusConflict, "User already exists")
		return
	}

	user := models.User{
		Username: registerRequest.Username,
		Name:     registerRequest.Name,
		Password: registerRequest.Password,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		utils.WriteErrorResponse(writer, http.StatusInternalServerError, "Failed to create a new user")
		return
	}

	utils.WriteSuccessResponse(writer, http.StatusCreated, "User registered successfully", nil)
}
