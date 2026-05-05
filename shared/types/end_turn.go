package types

type EndTurnRequest struct {
	PlayerId Uuid `json:"player_id"`
}

type EndTurnResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	JoinCode JoinCode `json:"join_code"`
}

func EndTurnError(message string) EndTurnResponse {
	return EndTurnResponse{
		Success: false,
		Message: message,
	}
}
