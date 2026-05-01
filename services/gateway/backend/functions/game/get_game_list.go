package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func GetGameList(api api.Api, req types.GetGameListRequest) types.GetGameListResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "get_game_list", req)

	if err != nil {
		return types.GetGameListError("")
	}

	var resp types.GetGameListResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.GetGameListError("")
	}

	return resp
}
