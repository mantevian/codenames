package handlers

import (
	"encoding/json"
	"net/http"

	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body     types.LoginRequest   true "Name and password"
// @Success      200         {object} types.LoginResponse "Login response with token if successful"
// @Router       /api/v1/login [post]
func Login(api Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.GenericResponse(w, http.StatusBadRequest, types.GenericResponseError("Invalid JSON"))
			return
		}

		responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "login", req)
		if err != nil {
			util.GenericResponse(w, http.StatusServiceUnavailable, types.GenericResponseError("Service unavailable"))
			return
		}

		var resp types.LoginResponse
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
