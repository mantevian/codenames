package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func JoinGame(api *api.Api, req types.JoinGameRequest) types.JoinGameResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "join_game", req)

	if err != nil {
		return types.JoinGameError("join_game: can't call")
	}

	var resp types.JoinGameResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.JoinGameError("join_game: response error")
	}

	return resp
}
