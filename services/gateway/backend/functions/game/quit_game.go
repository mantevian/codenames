package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func QuitGame(api *api.Api, req types.QuitGameRequest) types.QuitGameResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "quit_game", req)

	if err != nil {
		return types.QuitGameError("can't call")
	}

	var resp types.QuitGameResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.QuitGameError("response error")
	}

	return resp
}
