package types

import (
	"mantevian.xyz/codenames/shared/enums"
)

type CreateGameRequest struct {
	Language enums.Language `json:"language"`
}

type CreateGameResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Game    BasicGameResponse `json:"game"`
}

func CreateGameError(message string) CreateGameResponse {
	return CreateGameResponse{
		Success: false,
		Message: message,
	}
}
