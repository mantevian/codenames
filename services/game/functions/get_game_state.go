package functions

import (
	"database/sql"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
)

func GetGameState(req types.GetGameStateRequest, db *sql.DB) types.GetGameStateResponseMap {
	var err error

	var game types.Game
	var tiles []types.Tile
	var tilesHidden []types.Tile
	var turn types.Turn
	players := make(map[types.Uuid]types.Player) // user id!!!

	var rows *sql.Rows

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
			where join_code = $1
		`,
		req.JoinCode,
	)

	if err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "game not found")
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
			game_id,
			user_id,
			team,
			role,
			is_ready,
			created_at
		from players
		where game_id = $1
		`,
		game.Id,
	)

	if err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "player not found")
	}

	for rows.Next() {
		var player types.Player
		rows.Scan(
			&player.Id,
			&player.GameId,
			&player.UserId,
			&player.Team,
			&player.Role,
			&player.IsReady,
			&player.CreatedAt,
		)
		players[player.UserId] = player
	}
	rows.Close()

	rows, err = db.Query(`
		select
			position,
			game_id,
			type,
			is_revealed,
			word
		from tiles
		where game_id = $1
		order by position asc
		`,
		game.Id,
	)

	if err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "can't get game tiles")
	}

	for rows.Next() {
		var tile types.Tile

		rows.Scan(
			&tile.Position,
			&tile.GameId,
			&tile.Type,
			&tile.IsRevealed,
			&tile.Word,
		)

		tiles = append(tiles, tile)

		if tile.IsRevealed {
			tilesHidden = append(tilesHidden, tile)
		} else {
			tilesHidden = append(tilesHidden, types.Tile{
				Position:   tile.Position,
				GameId:     tile.GameId,
				IsRevealed: tile.IsRevealed,
				Word:       tile.Word,
			})
		}
	}
	rows.Close()

	rows, err = db.Query(`
		select
			player_id,
			clue_word,
			clue_number,
			guesses_left,
			created_at
		from turns
		where game_id = $1
		order by created_at desc
		limit 1
		`,
		game.Id,
	)

	if err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "turn not found")
	}

	if rows.Next() {
		rows.Scan(
			&turn.PlayerId,
			&turn.ClueWord,
			&turn.ClueNumber,
			&turn.GuessesLeft,
			&turn.CreatedAt,
		)
	}
	rows.Close()

	fullResults := make(types.GetGameStateResponseMap)

	for _, id := range req.UserIds {
		var res types.GetGameStateResponse

		switch players[id].Role {
		case enums.RoleOperative:
			res = types.GetGameStateResponse{
				Success: true,
				Game:    game,
				Tiles:   tilesHidden,
				Turn:    turn,
			}
		case enums.RoleSpymaster:
			res = types.GetGameStateResponse{
				Success: true,
				Game:    game,
				Tiles:   tiles,
				Turn:    turn,
			}
		}
		fullResults[id] = res
	}

	return fullResults
}
