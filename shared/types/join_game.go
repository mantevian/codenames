package types

type JoinGameRequest struct {
	UserId   Uuid   `json:"user_id"`
	JoinCode string `json:"join_code"`
}

type JoinGameResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func JoinGameError(message string) JoinGameResponse {
	return JoinGameResponse{
		Success: false,
		Message: message,
	}
}
