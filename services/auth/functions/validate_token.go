package functions

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_auth/jwt"
	"mantevian.xyz/codenames/shared/types"
)

func ValidateToken(payload []byte) types.ValidateTokenResponse {
	var req types.ValidateTokenRequest
	json.Unmarshal(payload, &req)
	claims, err := jwt.ValidateToken(req.Token)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	return types.ValidateTokenResponse{
		Success: true,
		Claims:  *claims,
	}
}
