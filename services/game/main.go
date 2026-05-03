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
	"mantevian.xyz/codenames/shared/types"
)

var db *sql.DB

func HandleRPC(action string, payload []byte) ([]byte, error) {
	var res any

	switch action {
	case "create_game":
		res = functions.CreateGame(payload, db)
	case "get_game_list":
		res = functions.GetGameList(db)
	case "join_game":
		res = functions.JoinGame(payload, db)
	case "quit_game":
		res = functions.QuitGame(payload, db)
	case "get_game_player_list":
		res = functions.GetGamePlayerList(payload, db)
	case "set_ready":
		res = functions.SetReady(payload, db)
	default:
		res = types.GenericResponseError(fmt.Sprintf("unknown rpc action %s", action))
	}

	return json.Marshal(res)
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
