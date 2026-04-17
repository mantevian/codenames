package functions

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_auth/jwt"
	"mantevian.xyz/codenames/shared/types"
)

func ValidateToken(payload []byte) types.ValidateTokenResponse {
	var token string
	json.Unmarshal(payload, &token)
	claims, err := jwt.ValidateToken(token)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	return types.ValidateTokenResponse{
		Success: true,
		Claims:  *claims,
	}
}
