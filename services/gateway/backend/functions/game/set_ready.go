package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func SetReady(api *api.Api, req types.SetReadyRequest) types.SetReadyResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "set_ready", req)

	if err != nil {
		return types.SetReadyError("")
	}

	var resp types.SetReadyResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.SetReadyError("")
	}

	return resp
}
