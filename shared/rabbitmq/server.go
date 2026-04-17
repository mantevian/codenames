package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"mantevian.xyz/codenames/shared/types"
)

type RPCServer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	handler func(action string, payload []byte) ([]byte, error)
}

func NewRPCServer(amqpURL, queue string) (*RPCServer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RPCServer{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (s *RPCServer) SetHandler(handler func(action string, payload []byte) ([]byte, error)) {
	s.handler = handler
}

func (s *RPCServer) Start() error {
	msgs, err := s.channel.Consume(
		s.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("RPC Server listening on queue: %s", s.queue)

	for msg := range msgs {
		var rpcMsg types.RPCMessage
		if err := json.Unmarshal(msg.Body, &rpcMsg); err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			continue
		}

		response, err := s.handler(rpcMsg.Action, rpcMsg.Payload)
		if err != nil {
			response = []byte(`{"error": "` + err.Error() + `"}`)
		}

		err = s.channel.Publish(
			"",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: msg.CorrelationId,
				Body:          response,
			},
		)
		if err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	}

	return nil
}

func (s *RPCServer) Close() {
	s.channel.Close()
	s.conn.Close()
}
