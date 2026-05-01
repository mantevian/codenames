package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func JoinGame(payload []byte, db *sql.DB) types.JoinGameResponse {
	var req types.JoinGameRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.JoinGameError(err.Error())
	}

	var gameId types.Uuid
	var userId types.Uuid = req.UserId
	var team enums.Team

	var redPlayers = 0
	var bluePlayers = 0

	rows, err := db.Query(`
			select
				games.id,
				players.team
			from games
			where join_code = $1
			inner join players on games.id = players.game_id
		`,
		req.JoinCode,
	)

	if err != nil {
		return types.JoinGameError("Can't find game")
	}

	for rows.Next() {
		var playerGameId types.Uuid
		var playerTeam enums.Team

		rows.Scan(&playerGameId, &playerTeam)

		gameId = playerGameId
		switch playerTeam {
		case enums.TeamRed:
			redPlayers++
		case enums.TeamBlue:
			bluePlayers++
		}
	}

	if redPlayers > bluePlayers {
		team = enums.TeamRed
	} else if redPlayers < bluePlayers {
		team = enums.TeamBlue
	} else {
		team = enums.RandomTeam()
	}

	rows, err = db.Query(`
		insert into players
			(
				id,
				game_id
				user_id,
				team,
				role,
				is_ready
			)
		values
			(
				gen_random_uuid(),
				$1,
				$2,
				$3,
				$4,
				false
			)
		returning *
		`,
		gameId,
		userId,
		team,
		enums.RoleOperative,
	)

	if err != nil {
		return types.JoinGameError("Cannot create player")
	}

	return types.JoinGameResponse{
		Success: true,
	}
}
