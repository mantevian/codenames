package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/service_auth/jwt"
	"mantevian.xyz/codenames/shared/types"
)

func ValidateToken(payload []byte, db *sql.DB) types.ValidateTokenResponse {
	var req types.ValidateTokenRequest
	json.Unmarshal(payload, &req)
	claims, err := jwt.ValidateToken(req.Token)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	res := types.ValidateTokenResponse{
		Success: true,
		Claims:  *claims,
	}

	userId := res.Claims.UserId

	rows, err := db.Query(
		"select id from users where id = $1",
		userId,
	)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	rows.Close()

	return res
}
