package util

import (
	"encoding/json"
	"net/http"

	"mantevian.xyz/codenames/shared/types"
)

func GenericResponse(w http.ResponseWriter, code int, resp types.GenericResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
