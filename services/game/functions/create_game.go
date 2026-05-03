package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func CreateGame(req types.CreateGameRequest, db *sql.DB) types.CreateGameResponse {
	rows, err := db.Query(`
		insert into games
			(
				id,
				starting_team,
				join_code,
				language,
				team_won,
				status,
				finished_at,
				created_at
			)
		values
			(
				gen_random_uuid(),
				$1,
				$2,
				$3,
				NULL,
				'waiting',
				NULL,
				now()
			)
		returning *
		`,
		enums.RandomTeam(),
		types.NewJoinCode(),
		req.Language,
	)

	if err != nil {
		return types.CreateGameError("Cannot create game")
	}

	var game types.BasicGameResponse

	rows.Next()
	rows.Scan(&game.Id, &game.StartingTeam, &game.JoinCode, &game.Language, nil, &game.Status, nil, &game.CreatedAt)

	return types.CreateGameResponse{
		Success: true,
		Game:    game,
	}
}
