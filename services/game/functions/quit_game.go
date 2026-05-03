package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/shared/types"
)

func QuitGame(payload []byte, db *sql.DB) types.QuitGameResponse {
	var req types.QuitGameRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.QuitGameError("can't parse request")
	}

	var gameId string
	var joinCode string

	rows, err := db.Query(`
			select
				games.id,
				games.join_code
			from games
			inner join players on players.game_id = games.id
			where players.id = $1
		`,
		req.PlayerId,
	)
	rows.Next()
	rows.Scan(&gameId, &joinCode)

	if err != nil {
		return types.QuitGameError("can't find game")
	}

	rows, err = db.Query(`
		delete from players
		where id = $1
		`,
		req.PlayerId,
	)

	if err != nil {
		return types.QuitGameError("cannot delete player")
	}

	return types.QuitGameResponse{
		Success:  true,
		JoinCode: joinCode,
	}
}
