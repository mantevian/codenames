package ws

import (
	"github.com/gorilla/websocket"
	"mantevian.xyz/codenames/shared/types"
)

type Client struct {
	Id     types.Uuid
	Conn   *websocket.Conn
	UserId types.Uuid
}
