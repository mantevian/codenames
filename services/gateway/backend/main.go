package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/service_gateway/gateway"
)

func main() {
	gateway, err := gateway.New()
	if err != nil {
		log.Fatal(err)
	}
	defer gateway.Close()

	var api = api.Api{Gateway: gateway}

	http.HandleFunc("POST /api/v1/register", api.Register)
	http.HandleFunc("POST /api/v1/login", api.Login)
	http.HandleFunc("POST /api/v1/validate_token", api.Auth(api.Ping))

	distPath := "../frontend/dist"
	fs := http.FileServer(http.Dir(distPath))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// serve existing files as is
		cleanPath := filepath.Clean(r.URL.Path)
		fullPath := filepath.Join(distPath, cleanPath)

		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// serve index.html on any non-matching path
		// then control over routing goes to the frontend
		r.URL.Path = "/"
		fs.ServeHTTP(w, r)
	})

	log.Printf("Gateway listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
