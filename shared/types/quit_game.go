package types

type QuitGameRequest struct {
	PlayerId Uuid `json:"player_id"`
}

type QuitGameResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	JoinCode JoinCode `json:"join_code"`
}

func QuitGameError(message string) QuitGameResponse {
	return QuitGameResponse{
		Success: false,
		Message: message,
	}
}
