package api

import (
	"encoding/json"
	"net/http"

	"mantevian.xyz/codenames/shared/types"
)

func (api *Api) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	responseBytes, err := api.Gateway.Call("login", req)
	if err != nil {
		http.Error(w, "Service unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	var resp types.LoginResponse
	if err := json.Unmarshal(responseBytes, &resp); err != nil {
		http.Error(w, "Invalid response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(resp)
}
