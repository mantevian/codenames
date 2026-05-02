package functions

import (
	"database/sql"
	"encoding/json"

	"mantevian.xyz/codenames/shared/types"
)

func GetGamePlayerList(payload []byte, db *sql.DB) types.GetGamePlayerListResponse {
	var req types.GetGamePlayerListRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return types.GetGamePlayerListError("can't parse request")
	}

	rows, err := db.Query(`
		select
			players.id,
			players.game_id,
			players.user_id,
			players.team,
			players.role,
			players.is_ready,
			users.name
		from players
		inner join games on players.game_id = games.id
		inner join users on players.user_id = users.id
		where games.join_code = $1
		`,
		req.JoinCode,
	)

	if err != nil {
		return types.GetGamePlayerListError("game not found")
	}

	var players []types.Player

	for rows.Next() {
		var player types.Player
		rows.Scan(&player.Id, &player.GameId, &player.UserId, &player.Team, &player.Role, &player.IsReady, &player.Name)
		players = append(players, player)
	}

	return types.GetGamePlayerListResponse{
		Success: true,
		Players: players,
	}
}
