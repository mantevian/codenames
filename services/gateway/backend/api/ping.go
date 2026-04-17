package api

import (
	"net/http"
)

func (api *Api) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
