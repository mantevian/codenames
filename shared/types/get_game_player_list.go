package types

type GetGamePlayerListRequest struct {
	JoinCode JoinCode `json:"join_code"`
}

type GetGamePlayerListResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Players []Player `json:"players"`
}

func GetGamePlayerListError(message string) GetGamePlayerListResponse {
	return GetGamePlayerListResponse{
		Success: false,
		Message: message,
	}
}
