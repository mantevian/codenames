package functions

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func Login(api api.Api, req types.LoginRequest) types.LoginResponse {
	responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "login", req)

	if err != nil {
		return types.LoginError("login: can't call")
	}

	var resp types.LoginResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.LoginError("login: call error")
	}

	return resp
}
