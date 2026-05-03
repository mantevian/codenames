package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func StartGame(api *api.Api, req types.StartGameRequest) types.StartGameResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "start_game", req)

	if err != nil {
		return types.StartGameError("")
	}

	var resp types.StartGameResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.StartGameError("")
	}

	return resp
}
