package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func JoinGame(req types.JoinGameRequest, db *sql.DB) types.JoinGameResponse {
	var gameId types.Uuid
	var gameStatus enums.GameStatus
	var userId types.Uuid = req.UserId
	var team enums.Team

	var redPlayers = 0
	var bluePlayers = 0

	rows, err := db.Query(`
			select
				id,
				status
			from games
			where
				join_code = $1
			and
				status <> 'finished'
			order by created_at desc
			limit 1
		`,
		req.JoinCode,
	)

	if err != nil {
		return types.JoinGameError("can't find game")
	}

	rows.Next()
	rows.Scan(&gameId, &gameStatus)
	rows.Close()

	if gameId == "" {
		return types.JoinGameError("can't find game")
	}

	rows, err = db.Query(`
			select
				user_id,
				team
			from players
			where game_id = $1
		`,
		gameId,
	)

	if err != nil {
		return types.JoinGameError("can't find game")
	}

	isPlayerInThisGame := false

	for rows.Next() {
		var playerUserId types.Uuid
		var playerTeam enums.Team

		rows.Scan(&playerUserId, &playerTeam)

		if playerUserId == req.UserId {
			isPlayerInThisGame = true
		}

		switch playerTeam {
		case enums.TeamRed:
			redPlayers++
		case enums.TeamBlue:
			bluePlayers++
		}
	}

	if !isPlayerInThisGame {
		if gameStatus == enums.GameStatusPlaying {
			return types.JoinGameError("this room is already playing")
		}
	}

	rows.Close()

	if redPlayers+bluePlayers == 4 {
		return types.JoinGameError("max 4 players")
	}

	if redPlayers > bluePlayers {
		team = enums.TeamBlue
	} else if redPlayers < bluePlayers {
		team = enums.TeamRed
	} else {
		team = enums.RandomTeam()
	}

	_, err = db.Exec(`
		insert into players
			(
				id,
				game_id,
				user_id,
				team,
				role,
				is_ready,
				created_at
			)
		values
			(
				gen_random_uuid(),
				$1,
				$2,
				$3,
				$4,
				false,
				now()
			)
		`,
		gameId,
		userId,
		team,
		enums.RoleOperative,
	)

	if err != nil {
		return types.JoinGameError("cannot create player")
	}

	return types.JoinGameResponse{
		Success: true,
	}
}
