package handlers

import (
	"encoding/json"
	"net/http"

	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func Register(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.GenericResponse(w, http.StatusBadRequest, types.GenericResponseError("Invalid JSON"))
			return
		}

		if req.Password != req.PasswordConfirm {
			util.GenericResponse(w, http.StatusBadRequest, types.GenericResponseError("Passwords don't match"))
			return
		}

		responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "register", req)
		if err != nil {
			util.GenericResponse(w, http.StatusServiceUnavailable, types.GenericResponseError("Service unavailable"))
			return
		}

		var resp types.RegisterResponse
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
