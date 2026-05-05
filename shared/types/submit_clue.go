package types

type SubmitClueRequest struct {
	PlayerId Uuid   `json:"player_id"`
	Word     string `json:"word"`
	Number   int    `json:"number"`
}

type SubmitClueResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	JoinCode JoinCode `json:"join_code"`
}

func SubmitClueError(message string) SubmitClueResponse {
	return SubmitClueResponse{
		Success: false,
		Message: message,
	}
}
