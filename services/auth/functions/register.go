package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/service_auth/hash"
	"mantevian.xyz/codenames/shared/types"
)

func Register(payload []byte, db *sql.DB) types.RegisterResponse {
	var req types.RegisterRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.RegisterError(err.Error())
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		return types.RegisterError(err.Error())
	}

	rows, err := db.Query(
		"insert into users (id, name, password, created_at) values (gen_random_uuid(), $1, $2, now()) returning *",
		req.Name,
		hash,
	)

	if err != nil {
		return types.RegisterError(err.Error())
	}

	var user types.User

	rows.Next()
	rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt)

	return types.RegisterResponse{
		Success:   true,
		UserId:    user.Id,
		Message:   "User registered successfully",
		CreatedAt: user.CreatedAt,
	}
}
