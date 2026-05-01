package handlers

import (
	"net/http"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/shared/types"
)

func Ping(api api.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.GenericResponse(w, http.StatusOK, types.GenericResponseSuccess("ok"))
	}
}
