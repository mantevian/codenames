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
		var req types.CreateGameRequest
		json.Unmarshal(payload, &req)
		res = functions.CreateGame(req, db)
	case "get_game_list":
		var req types.GetGameListRequest
		json.Unmarshal(payload, &req)
		res = functions.GetGameList(db)
	case "join_game":
		var req types.JoinGameRequest
		json.Unmarshal(payload, &req)
		res = functions.JoinGame(req, db)
	case "quit_game":
		var req types.QuitGameRequest
		json.Unmarshal(payload, &req)
		res = functions.QuitGame(req, db)
	case "get_game_player_list":
		var req types.GetGamePlayerListRequest
		json.Unmarshal(payload, &req)
		res = functions.GetGamePlayerList(req, db)
	case "set_ready":
		var req types.SetReadyRequest
		json.Unmarshal(payload, &req)
		res = functions.SetReady(req, db)
	case "start_game":
		var req types.StartGameRequest
		json.Unmarshal(payload, &req)
		res = functions.StartGame(req, db)
	case "get_game_state":
		var req types.GetGameStateRequest
		json.Unmarshal(payload, &req)
		res = functions.GetGameState(req, db)
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
