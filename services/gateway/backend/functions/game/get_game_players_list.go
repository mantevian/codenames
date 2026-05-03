package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func GetGamePlayersList(api *api.Api, req types.GetGamePlayerListRequest) types.GetGamePlayerListResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "get_game_player_list", req)
	if err != nil {
		return types.GetGamePlayerListError("")
	}

	var resp types.GetGamePlayerListResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.GetGamePlayerListError("")
	}

	return resp
}
