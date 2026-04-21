package game

import (
	"encoding/json"
	"net/http"

	"mantevian.xyz/codenames/service_gateway/handlers"
	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func GetWaitingGameList(api handlers.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWaitingGameListRequest

		responseBytes, err := api.Gateway.Call(rabbitmq.GameQueue, "get_waiting_game_list", req)
		if err != nil {
			util.GenericResponse(w, http.StatusServiceUnavailable, types.GenericResponseError("Service unavailable"))
			return
		}

		var resp types.GetWaitingGameListResponse
		if err := json.Unmarshal(responseBytes, &resp); err != nil {
			util.GenericResponse(w, http.StatusInternalServerError, types.GenericResponseError("Server error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if !resp.Success {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(resp)
	}
}
