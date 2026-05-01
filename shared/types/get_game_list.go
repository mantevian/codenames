package types

type GetGameListRequest struct {
}

type GetGameListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Games   []BasicGameResponse `json:"games"`
}

func GetGameListError(message string) GetGameListResponse {
	return GetGameListResponse{
		Success: false,
		Message: message,
	}
}
