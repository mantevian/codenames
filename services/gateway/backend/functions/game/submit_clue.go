package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func SubmitClue(api *api.Api, req types.SubmitClueRequest) types.SubmitClueResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "submit_clue", req)
	if err != nil {
		return types.SubmitClueError("can't call rpc method")
	}

	var resp types.SubmitClueResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.SubmitClueError("can't parse rpc response")
	}

	return resp
}
