package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func CreateGame(api api.Api, req types.CreateGameRequest) types.CreateGameResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "create_game", req)

	if err != nil {
		return types.CreateGameError("")
	}

	var resp types.CreateGameResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.CreateGameError("")
	}

	return resp
}
