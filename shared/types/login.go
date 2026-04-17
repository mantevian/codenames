package types

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func LoginError(message string) LoginResponse {
	return LoginResponse{
		Success: false,
		Message: message,
	}
}
