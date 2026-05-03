package functions

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func ValidateToken(api *api.Api, req types.ValidateTokenRequest) types.ValidateTokenResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "validate_token", req)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	var resp types.ValidateTokenResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.ValidateTokenFalse()
	}

	return resp
}
