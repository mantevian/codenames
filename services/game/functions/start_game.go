package functions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"mantevian.xyz/codenames/shared/enums"
	"mantevian.xyz/codenames/shared/types"
	"mantevian.xyz/codenames/shared/util"
)

func StartGame(payload []byte, db *sql.DB) types.StartGameResponse {
	var req types.StartGameRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.StartGameError("can't parse request")
	}

	var gameId string
	var startingTeam enums.Team
	var language enums.Language

	rows, err := db.Query(`
		select
			games.id,
			games.starting_team,
			games.language
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
	rows.Scan(&gameId, &startingTeam, &language)

	redTiles := 8
	blueTiles := 8

	switch startingTeam {
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
		language,
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

	var newRows []string

	for i := range 25 {
		word := words[i]
		tile := shuffledTiles[i]

		newRows = append(newRows,
			fmt.Sprintf(
				"(%d, '%s', '%s', %t, '%s')",
				i,
				gameId,
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
		gameId,
	)

	_, err = db.Exec(`
		update games
		set
			status = 'playing'
		where
			id = $1
		`,
		gameId,
	)

	if err != nil {
		return types.StartGameError("can't start game")
	}

	return types.StartGameResponse{
		Success: true,
	}
}
