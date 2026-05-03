package types

type SetReadyRequest struct {
	PlayerId Uuid `json:"player_id"`
	IsReady  bool `json:"is_ready"`
}

type SetReadyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SetReadyError(message string) SetReadyResponse {
	return SetReadyResponse{
		Success: false,
		Message: message,
	}
}
