package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func EndTurn(api *api.Api, req types.EndTurnRequest) types.EndTurnResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "end_turn", req)
	if err != nil {
		return types.EndTurnError("can't call rpc method")
	}

	var resp types.EndTurnResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.EndTurnError("can't parse rpc response")
	}

	return resp
}
