package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/types"
)

func QuitGame(req types.QuitGameRequest, db *sql.DB) types.QuitGameResponse {
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
