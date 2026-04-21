package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"mantevian.xyz/codenames/service_game/functions"
	"mantevian.xyz/codenames/shared/rabbitmq"
)

var db *sql.DB

func HandleRPC(action string, payload []byte) ([]byte, error) {
	switch action {
	case "create_game":
		res := functions.CreateGame(payload, db)
		return json.Marshal(res)
	case "get_waiting_game_list":
		res := functions.GetWaitingGameList(db)
		return json.Marshal(res)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server, err := rabbitmq.NewRPCServer(os.Getenv("RABBITMQ_URL"), rabbitmq.GameQueue)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	server.SetHandler(HandleRPC)

	log.Println("Game service starting...")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
