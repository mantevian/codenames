package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"mantevian.xyz/codenames/service_gateway/api"
	"mantevian.xyz/codenames/service_gateway/functions"
	"mantevian.xyz/codenames/service_gateway/functions/game"
	"mantevian.xyz/codenames/service_gateway/util"
	"mantevian.xyz/codenames/service_gateway/ws"
	"mantevian.xyz/codenames/shared/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpdatePlayerListByJoinCode(api *api.Api, hub *ws.Hub, joinCode string) {
	res := game.GetGamePlayersList(
		api,
		types.GetGamePlayerListRequest{
			JoinCode: joinCode,
		},
	)
	json, _ := json.Marshal(res)
	hub.Broadcast(joinCode, "update_player_list", json)
}

func UpdatePlayerListByPlayerId(api *api.Api, hub *ws.Hub, playerId string) {
	res := game.GetGamePlayersList(
		api,
		types.GetGamePlayerListRequest{
			PlayerId: playerId,
		},
	)
	json, _ := json.Marshal(res)
	hub.Broadcast(res.JoinCode, "update_player_list", json)
}

func UpdateGameStateByJoinCode(api *api.Api, hub *ws.Hub, joinCode string) {
	clients := hub.GetClientsInGroup(joinCode)
	var userIds []types.Uuid

	for _, c := range clients {
		userIds = append(userIds, c.UserId)
	}

	res := game.GetGameState(
		api,
		types.GetGameStateRequest{
			JoinCode: joinCode,
			UserIds:  userIds,
		},
	)

	for _, c := range clients {
		resJson, _ := json.Marshal(res[c.UserId])
		message := types.WsMessage{
			Action:  "update_game_state",
			Payload: resJson,
		}

		messageJson, _ := json.Marshal(message)
		c.Conn.WriteMessage(websocket.TextMessage, messageJson)

		fmt.Printf("send >> %s >> %s\n", c.Id, string(resJson))
	}
}

func Ws(api *api.Api, hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			util.GenericResponse(w, http.StatusInternalServerError, types.GenericResponseError("Server error"))
			return
		}

		id := uuid.New().String()

		client := ws.Client{
			Id:   id,
			Conn: conn,
		}

		hub.AddClient("", &client)

		defer conn.Close()
		defer hub.RemoveClient(id)

		for {
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println("incoming >>", string(bytes))

			var message types.WsMessage
			json.Unmarshal(bytes, &message)

			validateTokenRes := functions.ValidateToken(api, types.ValidateTokenRequest{Token: message.Token})
			authSuccess := validateTokenRes.Success

			hub.SetClientUserId(client.Id, validateTokenRes.Claims.UserId)

			var res any

			switch message.Action {
			case "ping":
				res = functions.Ping(api)
			case "register":
				var req types.RegisterRequest
				json.Unmarshal(message.Payload, &req)
				res = functions.Register(api, req)
			case "login":
				var req types.LoginRequest
				json.Unmarshal(message.Payload, &req)
				loginRes := functions.Login(api, req)
				res = loginRes
				hub.SetClientUserId(client.Id, loginRes.UserId)
			case "validate_token":
				res = validateTokenRes
			}

			if authSuccess {
				userId := validateTokenRes.Claims.UserId

				switch message.Action {
				case "create_game":
					var req types.CreateGameRequest
					json.Unmarshal(message.Payload, &req)
					res = game.CreateGame(api, req)
				case "get_game_list":
					var req types.GetGameListRequest
					json.Unmarshal(message.Payload, &req)
					res = game.GetGameList(api, req)
				case "join_game":
					var req types.JoinGameRequest
					json.Unmarshal(message.Payload, &req)
					req.UserId = userId
					joinGameRes := game.JoinGame(api, req)
					res = joinGameRes

					if joinGameRes.Success {
						hub.MoveClient(id, req.JoinCode)
						UpdatePlayerListByJoinCode(api, hub, req.JoinCode)
						UpdateGameStateByJoinCode(api, hub, req.JoinCode)
					}
				case "quit_game":
					var req types.QuitGameRequest
					json.Unmarshal(message.Payload, &req)
					quitGameRes := game.QuitGame(api, req)
					res = quitGameRes

					if quitGameRes.Success {
						hub.MoveClient(id, "")
						UpdatePlayerListByJoinCode(api, hub, quitGameRes.JoinCode)
					}
				case "set_ready":
					var req types.SetReadyRequest
					json.Unmarshal(message.Payload, &req)
					setReadyRes := game.SetReady(api, req)
					res = setReadyRes

					if setReadyRes.Success {
						UpdatePlayerListByPlayerId(api, hub, req.PlayerId)
					}
				case "start_game":
					var req types.StartGameRequest
					json.Unmarshal(message.Payload, &req)
					startGameRes := game.StartGame(api, req)
					res = startGameRes

					if startGameRes.Success {
						UpdatePlayerListByPlayerId(api, hub, req.PlayerId)
						hub.Broadcast(startGameRes.Game.JoinCode, "game_started", []byte("{}"))
						UpdateGameStateByJoinCode(api, hub, startGameRes.Game.JoinCode)
					}
				case "submit_clue":
					var req types.SubmitClueRequest
					json.Unmarshal(message.Payload, &req)
					submitClueRes := game.SubmitClue(api, req)
					res = submitClueRes

					if submitClueRes.Success {
						UpdateGameStateByJoinCode(api, hub, submitClueRes.JoinCode)
					}
				case "submit_guess":
					var req types.SubmitGuessRequest
					json.Unmarshal(message.Payload, &req)
					submitGuessRes := game.SubmitGuess(api, req)
					res = submitGuessRes

					if submitGuessRes.Success {
						UpdateGameStateByJoinCode(api, hub, submitGuessRes.JoinCode)
					}
				case "end_turn":
					var req types.EndTurnRequest
					json.Unmarshal(message.Payload, &req)
					endTurnRes := game.EndTurn(api, req)
					res = endTurnRes

					if endTurnRes.Success {
						UpdateGameStateByJoinCode(api, hub, endTurnRes.JoinCode)
					}
				}
			}

			response := types.WsResponse{
				Action:  message.Action,
				Id:      message.Id,
				Payload: res,
			}

			responseJson, _ := json.Marshal(response)
			fmt.Println("response >>", string(responseJson))

			if err := conn.WriteMessage(websocket.TextMessage, responseJson); err != nil {
				break
			}
		}
	}
}
