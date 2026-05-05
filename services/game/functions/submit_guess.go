package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func SubmitGuess(req types.SubmitGuessRequest, db *sql.DB) types.SubmitGuessResponse {
	var err error

	var game types.Game
	var player types.Player
	var turn types.Turn
	var tile types.Tile
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
		return types.SubmitGuessError("player not found")
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
		return types.SubmitGuessError("game not found")
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
				guesses_left,
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
		return types.SubmitGuessError("turn not found")
	}

	rows.Next()
	rows.Scan(
		&turn.Id,
		&turn.PlayerId,
		&turn.GameId,
		&turn.ClueWord,
		&turn.ClueNumber,
		&turn.GuessesLeft,
		&turn.CreatedAt,
	)
	rows.Close()

	isMyRole := player.Role == game.CurrentTurnRole
	isMyTeam := player.Team == game.CurrentTurnTeam
	areGuessesLeft := turn.GuessesLeft > 0
	canMakeTurn := isMyRole && isMyTeam && areGuessesLeft

	if !canMakeTurn {
		return types.SubmitGuessError("can't make a turn")
	}

	if req.Position < 0 || req.Position > 24 {
		return types.SubmitGuessError("position should be in range 0..24")
	}

	rows, err = db.Query(`
			select
				type,
				is_revealed
			from tiles
			where
				game_id = $1
			and
				position = $2
		`,
		game.Id,
		req.Position,
	)

	if err != nil {
		return types.SubmitGuessError("tile not found")
	}

	rows.Next()
	rows.Scan(
		&tile.Type,
		&tile.IsRevealed,
	)
	rows.Close()

	if tile.IsRevealed {
		return types.SubmitGuessError("this tile is already revealed")
	}

	_, err = db.Exec(`
		insert into guesses
			(
				turn_id,
				position,
				created_at
			)
		values
			(
				$1,
				$2,
				now()
			)
		`,
		turn.Id,
		req.Position,
	)

	if err != nil {
		return types.SubmitGuessError("cannot submit guess")
	}

	if tile.Type == enums.TileNeutral {
		turn.GuessesLeft = 0
	} else {
		turn.GuessesLeft -= 1
	}

	_, err = db.Exec(`
			update tiles
			set
				is_revealed = true
			where
				game_id = $1
			and
				position = $2
		`,
		game.Id,
		req.Position,
	)

	if err != nil {
		return types.SubmitGuessError("cannot update tile")
	}

	_, err = db.Exec(`
			update turns
			set
				guesses_left = $2
			where
				id = $1
		`,
		turn.Id,
		turn.GuessesLeft,
	)

	if err != nil {
		return types.SubmitGuessError("cannot update tile")
	}

	var countRed int
	var countBlue int

	rows, err = db.Query(`
			select
				count(*) filter (where type = 'red' and is_revealed = true) as red_count,
				count(*) filter (where type = 'blue' and is_revealed = true) as blue_count
			from tiles
			where
				game_id = $1
		`,
		game.Id,
	)

	if err != nil {
		return types.SubmitGuessError("tiles not found")
	}

	rows.Next()
	rows.Scan(
		&countRed,
		&countBlue,
	)
	rows.Close()

	var teamWon enums.Team

	if tile.Type == enums.TileAssassin {
		if player.Team == enums.TeamRed {
			teamWon = enums.TeamBlue
		} else {
			teamWon = enums.TeamRed
		}
	}

	if countRed == 9 {
		teamWon = enums.TeamRed
	}

	if countBlue == 9 {
		teamWon = enums.TeamBlue
	}

	println(countRed, countBlue, teamWon)

	if teamWon != "" {
		_, err = db.Exec(`
			update games
			set
				status = 'finished',
				team_won = $2,
				finished_at = now()
			where
				id = $1
		`,
			game.Id,
			teamWon,
		)

		if err != nil {
			return types.SubmitGuessError("cannot update tile")
		}
	}

	return types.SubmitGuessResponse{
		Success:  true,
		Message:  "guess submitted!",
		JoinCode: game.JoinCode,
	}
}
