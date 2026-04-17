package main

import (
	"log"
	"net/http"

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

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	log.Printf("Gateway listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
