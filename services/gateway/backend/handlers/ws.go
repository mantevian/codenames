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

func Ws(api api.Api, hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			util.GenericResponse(w, http.StatusInternalServerError, types.GenericResponseError("Server error"))
			return
		}

		id := uuid.New().String()

		hub.AddClient("", &ws.Client{
			Id:   id,
			Conn: conn,
		})

		defer conn.Close()
		defer hub.RemoveClient(id)

		for {
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println(">>>", string(bytes))

			var message types.WsMessage
			json.Unmarshal(bytes, &message)

			validateTokenRes := functions.ValidateToken(api, types.ValidateTokenRequest{Token: message.Token})
			authSuccess := validateTokenRes.Success

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
				res = functions.Login(api, req)
			case "validate_token":
				res = validateTokenRes
			default:
				res = types.GenericResponseError("not logged in")
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
						playerListRes := game.GetGamePlayersList(
							api,
							types.GetGamePlayerListRequest{
								JoinCode: req.JoinCode,
							},
						)
						json, _ := json.Marshal(playerListRes)
						hub.Broadcast(req.JoinCode, "update_player_list", json)
					}
				}
			}

			response := types.WsResponse{
				Action:  message.Action,
				Id:      message.Id,
				Payload: res,
			}

			responseJson, _ := json.Marshal(response)
			fmt.Println("response:", string(responseJson))

			if err := conn.WriteMessage(websocket.TextMessage, responseJson); err != nil {
				break
			}
		}
	}
}
