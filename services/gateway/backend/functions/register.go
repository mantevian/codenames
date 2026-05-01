package functions

import (
	"encoding/json"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func Register(api api.Api, req types.RegisterRequest) types.RegisterResponse {
	if req.Password != req.PasswordConfirm {
		return types.RegisterError("Passwords don't match")
	}

	responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "register", req)

	if err != nil {
		return types.RegisterError("")
	}

	var resp types.RegisterResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		return types.RegisterError("")
	}

	return resp
}
