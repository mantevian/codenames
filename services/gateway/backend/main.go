package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "mantevian.xyz/codenames/service_gateway/docs"
	"mantevian.xyz/codenames/service_gateway/handlers/game"
	"mantevian.xyz/codenames/service_gateway/middleware"

	"mantevian.xyz/codenames/service_gateway/gateway"
	"mantevian.xyz/codenames/service_gateway/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	gateway, err := gateway.New()
	if err != nil {
		log.Fatal(err)
	}
	defer gateway.Close()

	var api = handlers.Api{Gateway: gateway}

	http.HandleFunc("POST /api/v1/register", handlers.Register(api))
	http.HandleFunc("POST /api/v1/login", handlers.Login(api))
	http.HandleFunc("POST /api/v1/validate_token", middleware.Auth(api)(handlers.Ping(api)))
	http.HandleFunc("POST /api/v1/create_game", middleware.Auth(api)(game.CreateGame(api)))
	http.HandleFunc("POST /api/v1/get_waiting_game_list", middleware.Auth(api)(game.GetWaitingGameList(api)))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

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
