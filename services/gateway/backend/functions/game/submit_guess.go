package game

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func SubmitGuess(api *api.Api, req types.SubmitGuessRequest) types.SubmitGuessResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "submit_guess", req)
	if err != nil {
		return types.SubmitGuessError("can't call rpc method")
	}

	var resp types.SubmitGuessResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.SubmitGuessError("can't parse rpc response")
	}

	return resp
}
