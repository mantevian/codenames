package types

type StartGameRequest struct {
	PlayerId Uuid `json:"player_id"`
}

type StartGameResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Game    Game   `json:"game"`
}

func StartGameError(message string) StartGameResponse {
	return StartGameResponse{
		Success: false,
		Message: message,
	}
}
