package types

type GetGamePlayerListRequest struct {
	JoinCode JoinCode `json:"join_code"`
	PlayerId Uuid     `json:"player_id"`
}

type GetGamePlayerListResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	JoinCode JoinCode `json:"join_code"`
	Players  []Player `json:"players"`
}

func GetGamePlayerListError(message string) GetGamePlayerListResponse {
	return GetGamePlayerListResponse{
		Success: false,
		Message: message,
	}
}
