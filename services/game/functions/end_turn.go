package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func EndTurn(req types.EndTurnRequest, db *sql.DB) types.EndTurnResponse {
	var err error

	var game types.Game
	var player types.Player
	var turn types.Turn
	var rows *sql.Rows
	var nextTeam enums.Team

	rows, err = db.Query(`
		select
			id,
			game_id,
			user_id,
			team,
			role,
			is_ready,
			created_at
		from players
		where id = $1
		`,
		req.PlayerId,
	)

	if err != nil {
		return types.EndTurnError("player not found")
	}

	rows.Next()
	rows.Scan(
		&player.Id,
		&player.GameId,
		&player.UserId,
		&player.Team,
		&player.Role,
		&player.IsReady,
		&player.CreatedAt,
	)

	rows.Close()

	rows, err = db.Query(`
			select
				id,
				starting_team,
				current_turn_team,
				current_turn_role,
				join_code,
				language,
				team_won,
				status,
				finished_at,
				created_at
			from games
			where
				id = $1
			and
				status = 'playing'
		`,
		player.GameId,
	)

	if err != nil {
		return types.EndTurnError("game not found")
	}

	rows.Next()
	rows.Scan(
		&game.Id,
		&game.StartingTeam,
		&game.CurrentTurnTeam,
		&game.CurrentTurnRole,
		&game.JoinCode,
		&game.Language,
		&game.TeamWon,
		&game.Status,
		&game.FinishedAt,
		&game.CreatedAt,
	)
	rows.Close()

	rows, err = db.Query(`
			select
				id,
				player_id,
				game_id,
				clue_word,
				clue_number,
				created_at
			from turns
			where
				game_id = $1
			order by created_at desc
			limit 1
		`,
		game.Id,
	)

	if err != nil {
		return types.EndTurnError("turn not found")
	}

	rows.Next()
	rows.Scan(
		&turn.PlayerId,
		&turn.GameId,
		&turn.ClueWord,
		&turn.ClueNumber,
		&turn.CreatedAt,
	)
	rows.Close()

	canMakeTurn := player.Role == game.CurrentTurnRole && player.Team == game.CurrentTurnTeam

	if !canMakeTurn {
		return types.EndTurnError("not your turn")
	}

	if game.CurrentTurnTeam == enums.TeamRed {
		nextTeam = enums.TeamBlue
	} else {
		nextTeam = enums.TeamRed
	}

	_, err = db.Exec(`
		update games
		set
			current_turn_role = 'spymaster',
			current_turn_team = $2
		where
			id = $1
		`,
		game.Id,
		nextTeam,
	)

	if err != nil {
		return types.EndTurnError("cannot set game state")
	}

	return types.EndTurnResponse{
		Success:  true,
		Message:  "turn ended",
		JoinCode: game.JoinCode,
	}
}
