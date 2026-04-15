package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

// MockDB simulates database operations
type MockDB struct {
	users map[string]types.RegisterRequest
}

func NewMockDB() *MockDB {
	return &MockDB{
		users: make(map[string]types.RegisterRequest),
	}
}

func (db *MockDB) CreateUser(req types.RegisterRequest) (string, error) {
	// Simulate validation
	if req.Email == "" || req.Password == "" {
		return "", fmt.Errorf("email and password required")
	}

	// Check if user exists
	if _, exists := db.users[req.Email]; exists {
		return "", fmt.Errorf("user already exists")
	}

	// Create user
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	db.users[req.Email] = req

	return userID, nil
}

type AuthService struct {
	db *MockDB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: NewMockDB(),
	}
}

func (a *AuthService) HandleRPC(action string, payload []byte) ([]byte, error) {
	switch action {
	case "register":
		return a.handleRegister(payload)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (a *AuthService) handleRegister(payload []byte) ([]byte, error) {
	var req types.RegisterRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		response := types.RegisterResponse{
			Success: false,
			Message: "Invalid request format",
		}
		return json.Marshal(response)
	}

	userID, err := a.db.CreateUser(req)
	if err != nil {
		response := types.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}
		return json.Marshal(response)
	}

	response := types.RegisterResponse{
		Success:   true,
		UserID:    userID,
		Message:   "User registered successfully",
		CreatedAt: time.Now(),
	}
	return json.Marshal(response)
}

func main() {
	authService := NewAuthService()

	server, err := rabbitmq.NewRPCServer(os.Getenv("RABBITMQ_URL"), rabbitmq.AuthQueue)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	server.SetHandler(authService.HandleRPC)

	log.Println("Auth service starting...")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
