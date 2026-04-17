package gateway

import (
	"context"
	"os"
	"time"

	"mantevian.xyz/codenames/shared/rabbitmq"
)

type Gateway struct {
	rpcClient *rabbitmq.RPCClient
}

func New() (*Gateway, error) {
	client, err := rabbitmq.NewRPCClient(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}
	return &Gateway{rpcClient: client}, nil
}

func (g *Gateway) Close() {
	g.rpcClient.Close()
}

func (g *Gateway) Call(action string, payload any) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	responseBytes, err := g.rpcClient.Call(ctx, rabbitmq.AuthQueue, action, payload)

	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}
