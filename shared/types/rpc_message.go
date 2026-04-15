package types

import "encoding/json"

type RPCMessage struct {
	Action        string          `json:"action"`
	Payload       json.RawMessage `json:"payload"`
	CorrelationID string          `json:"correlation_id"`
	ReplyTo       string          `json:"reply_to"`
}
