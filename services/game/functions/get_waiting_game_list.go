package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func GetWaitingGameList(db *sql.DB) types.GetWaitingGameListResponse {
	rows, err := db.Query(`
		select
			id,
			starting_team,
			join_code,
			language,
			created_at
		from games
		where status = 'waiting'
		`,
	)

	if err != nil {
		return types.GetWaitingGameListError("no games :(")
	}

	var games []types.BasicGameResponse

	for rows.Next() {
		var game types.BasicGameResponse
		rows.Scan(&game.Id, &game.StartingTeam, &game.JoinCode, &game.Language, &game.CreatedAt)
		game.Status = enums.GameStatusWaiting
		games = append(games, game)
	}

	return types.GetWaitingGameListResponse{
		Success: true,
		Games:   games,
	}
}
