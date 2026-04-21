package types

type GetWaitingGameListRequest struct {
}

type GetWaitingGameListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Games   []BasicGameResponse `json:"games"`
}

func GetWaitingGameListError(message string) GetWaitingGameListResponse {
	return GetWaitingGameListResponse{
		Success: false,
		Message: message,
	}
}
