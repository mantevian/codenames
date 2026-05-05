package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/types"
)

func SubmitClue(req types.SubmitClueRequest, db *sql.DB) types.SubmitClueResponse {
	var err error

	var game types.Game
	var player types.Player
	var rows *sql.Rows

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
		return types.SubmitClueError("player not found")
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
		return types.SubmitClueError("game not found")
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

	canMakeTurn := player.Role == game.CurrentTurnRole && player.Team == game.CurrentTurnTeam

	if !canMakeTurn {
		return types.SubmitClueError("not your turn")
	}

	if len(req.Word) < 1 || len(req.Word) > 32 {
		return types.SubmitClueError("word length should be in range 1..32")
	}

	if req.Number < 1 || req.Number > 4 {
		return types.SubmitClueError("number should be in range 1..4")
	}

	_, err = db.Exec(`
		insert into turns
			(
				id,
				player_id,
				game_id,
				clue_word,
				clue_number,
				guesses_left,
				created_at
			)
		values
			(
				gen_random_uuid(),
				$1,
				$2,
				$3,
				$4,
				2,
				now()
			)
		`,
		player.Id,
		game.Id,
		req.Word,
		req.Number,
	)

	if err != nil {
		return types.SubmitClueError("cannot submit clue")
	}

	_, err = db.Exec(`
		update games
		set
			current_turn_role = 'operative'
		where
			id = $1
		`,
		game.Id,
	)

	if err != nil {
		return types.SubmitClueError("cannot set game state")
	}

	return types.SubmitClueResponse{
		Success:  true,
		Message:  "clue submitted!",
		JoinCode: game.JoinCode,
	}
}
