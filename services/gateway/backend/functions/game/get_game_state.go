package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func GetGameState(api *api.Api, req types.GetGameStateRequest) types.GetGameStateResponseMap {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "get_game_state", req)
	if err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "can't call rpc method")
	}

	var resp types.GetGameStateResponseMap
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.GetGameStateErrorMap(req.UserIds, "can't parse rpc response")
	}

	return resp
}
