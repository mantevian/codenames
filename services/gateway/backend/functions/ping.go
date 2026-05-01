package functions

import (
	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/shared/types"
)

func Ping(api api.Api) types.GenericResponse {
	return types.GenericResponseSuccess("ok")
}
