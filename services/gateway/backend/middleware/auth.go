package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"mantevian.xyz/codenames/service_gateway/handlers"
	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

func Auth(api handlers.Api) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				util.GenericResponse(w, http.StatusUnauthorized, types.GenericResponseError("Missing token"))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			responseBytes, err := api.Gateway.Call(rabbitmq.AuthQueue, "validate_token", tokenString)

			var resp types.ValidateTokenResponse
			if err := json.Unmarshal(responseBytes, &resp); err != nil {
				util.GenericResponse(w, http.StatusInternalServerError, types.GenericResponseError("Server error"))
				return
			}

			if !resp.Success || err != nil {
				util.GenericResponse(w, http.StatusUnauthorized, types.GenericResponseError("Invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), "user", resp.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
