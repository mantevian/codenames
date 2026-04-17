package functions

import (
	"mantevian.xyz/codenames/service_auth/jwt"
	"mantevian.xyz/codenames/shared/types"
)

func ValidateToken(payload []byte) types.ValidateTokenResponse {
	token := string(payload)
	claims, err := jwt.ValidateToken(token)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	return types.ValidateTokenResponse{
		Success: true,
		Claims:  *claims,
	}
}
