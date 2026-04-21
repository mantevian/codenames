package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/service_auth/hash"
	"mantevian.xyz/codenames/service_auth/jwt"
	"mantevian.xyz/codenames/shared/types"
)

func Login(payload []byte, db *sql.DB) types.LoginResponse {
	var req types.LoginRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.LoginError(err.Error())
	}

	rows, err := db.Query(
		"select id, name, password, created_at from users where name = $1",
		req.Name,
	)

	if err != nil {
		return types.LoginError(err.Error())
	}

	var user types.User

	rows.Next()
	rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt)

	if !hash.CheckPassword(req.Password, user.Password) {
		return types.LoginError("Incorrect credentials")
	}

	token, err := jwt.GenerateToken(user.Id, user.Name)
	if err != nil {
		return types.LoginError(err.Error())
	}

	return types.LoginResponse{
		Success: true,
		UserId:  user.Id,
		Message: "User logged in successfully",
		Token:   token,
	}
}
