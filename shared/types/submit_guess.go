package types

type SubmitGuessRequest struct {
	PlayerId Uuid `json:"player_id"`
	Position int  `json:"position"`
}

type SubmitGuessResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	JoinCode JoinCode `json:"join_code"`
}

func SubmitGuessError(message string) SubmitGuessResponse {
	return SubmitGuessResponse{
		Success: false,
		Message: message,
	}
}
