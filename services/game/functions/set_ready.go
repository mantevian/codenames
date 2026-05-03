package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/shared/types"
)

func SetReady(payload []byte, db *sql.DB) types.SetReadyResponse {
	var req types.SetReadyRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.SetReadyError("can't parse request")
	}

	_, err = db.Exec(`
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
