package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"mantevian.xyz/codenames/shared/rabbitmq"
	"mantevian.xyz/codenames/shared/types"
)

type Gateway struct {
	rpcClient *rabbitmq.RPCClient
}

func NewGateway() (*Gateway, error) {
	client, err := rabbitmq.NewRPCClient(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}
	return &Gateway{rpcClient: client}, nil
}

func (g *Gateway) Close() {
	g.rpcClient.Close()
}

func (g *Gateway) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Call auth service via RPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	responseBytes, err := g.rpcClient.Call(ctx, rabbitmq.AuthQueue, "register", req)
	if err != nil {
		http.Error(w, "Service unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	var resp types.RegisterResponse
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

func main() {
	gateway, err := NewGateway()
	if err != nil {
		log.Fatal(err)
	}
	defer gateway.Close()

	http.HandleFunc("/register", gateway.RegisterHandler)

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	log.Printf("Gateway listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
