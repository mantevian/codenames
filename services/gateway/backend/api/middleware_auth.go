package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"mantevian.xyz/codenames/shared/types"
)

func (api *Api) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		responseBytes, err := api.Gateway.Call("validate_token", tokenString)

		var resp types.ValidateTokenResponse
		if err := json.Unmarshal(responseBytes, &resp); err != nil {
			http.Error(w, "Invalid response", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", resp.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
