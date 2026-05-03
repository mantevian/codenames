package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/types"
)

func SetReady(req types.SetReadyRequest, db *sql.DB) types.SetReadyResponse {
	_, err := db.Exec(`
		update players
		set is_ready = $2
		where id = $1
		`,
		req.PlayerId,
		req.IsReady,
	)

	if err != nil {
		return types.SetReadyError("cannot set ready")
	}

	return types.SetReadyResponse{
		Success: true,
	}
}
