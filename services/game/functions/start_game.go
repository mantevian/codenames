package functions

import (
	"database/sql"
	"fmt"
	"strings"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
	"mantevian.xyz/codenames/shared/util"
)

func StartGame(req types.StartGameRequest, db *sql.DB) types.StartGameResponse {
	var game types.Game
	var playersReady []bool

	rows, err := db.Query(`
		select
			games.id,
			games.starting_team,
			games.join_code,
			games.language,
			games.status,
			games.created_at
		from games
		inner join players on players.game_id = games.id
		where
			players.id = $1
		`,
		req.PlayerId,
	)

	if err != nil {
		return types.StartGameError("can't find game")
	}

	rows.Next()
	rows.Scan(&game.Id, &game.StartingTeam, &game.JoinCode, &game.Language, &game.Status, &game.CreatedAt)
	rows.Close()

	if game.Status != enums.GameStatusWaiting {
		return types.StartGameError("game already started")
	}

	rows, err = db.Query(`
		select
			is_ready
		from players
		where game_id = $1
		`,
		game.Id,
	)

	if err != nil {
		return types.StartGameError("can't find players in game")
	}

	for rows.Next() {
		var ready bool
		rows.Scan(&ready)
		playersReady = append(playersReady, ready)

		if !ready {
			return types.StartGameError("not all players are ready to start")
		}
	}

	rows.Close()

	if len(playersReady) != 4 {
		return types.StartGameError("need 4 players to start")
	}

	redTiles := 8
	blueTiles := 8

	switch game.StartingTeam {
	case enums.TeamRed:
		redTiles = 9
	case enums.TeamBlue:
		blueTiles = 9
	}

	var shuffledTiles = util.MakeShuffledTileList(redTiles, blueTiles, 7, 1)

	rows, err = db.Query(`
		select
			word
		from words
		where language = $1
		order by random()
		limit 25
		`,
		game.Language,
	)

	if err != nil {
		return types.StartGameError("can't get words")
	}

	var words []string

	for rows.Next() {
		var word string
		rows.Scan(&word)
		words = append(words, word)
	}

	rows.Close()

	var newRows []string

	for i := range 25 {
		word := words[i]
		tile := shuffledTiles[i]

		newRows = append(newRows,
			fmt.Sprintf(
				"(%d, '%s', '%s', %t, '%s')",
				i,
				game.Id,
				tile,
				false,
				word,
			),
		)
	}

	_, err = db.Exec(fmt.Sprintf(`
		insert into tiles
			(position, game_id, type, is_revealed, word)
		values
			%s
	`, strings.Join(newRows, ",")))

	if err != nil {
		return types.StartGameError("can't create tiles")
	}

	_, err = db.Exec(`
		with chosen as (
			select id
			from (
				select
					id,
					row_number() over (partition by team order by random()) as rn
				from players
				where game_id = $1
			) t
			where rn = 1
		)
		update players
		set role = 'spymaster'
		where id in (select id from chosen)
	`,
		game.Id,
	)

	_, err = db.Exec(`
		update games
		set
			status = 'playing',
			current_turn_team = $2,
			current_turn_role = 'spymaster'
		where
			id = $1
		`,
		game.Id,
		game.StartingTeam,
	)

	game.CurrentTurnTeam = game.StartingTeam
	game.CurrentTurnRole = enums.RoleSpymaster

	if err != nil {
		return types.StartGameError("can't start game")
	}

	return types.StartGameResponse{
		Success: true,
		Game:    game,
	}
}
