package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"mantevian.xyz/codenames/shared/types"
)

type RPCClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	pending map[string]chan []byte
	mu      sync.RWMutex
}

func NewRPCClient(amqpURL string) (*RPCClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Single consumer for direct reply-to
	msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	client := &RPCClient{
		conn:    conn,
		channel: ch,
		msgs:    msgs,
		pending: make(map[string]chan []byte),
	}

	// Start response dispatcher
	go client.dispatcher()

	return client, nil
}

func (c *RPCClient) dispatcher() {
	for msg := range c.msgs {
		c.mu.RLock()
		ch, ok := c.pending[msg.CorrelationId]
		c.mu.RUnlock()

		if ok {
			ch <- msg.Body
		}
	}
}

func (c *RPCClient) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *RPCClient) Call(ctx context.Context, queue string, action string, payload any) ([]byte, error) {
	corrID := fmt.Sprintf("%d", time.Now().UnixNano())
	respChan := make(chan []byte, 1)

	c.mu.Lock()
	c.pending[corrID] = respChan
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, corrID)
		c.mu.Unlock()
	}()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	msg := types.RPCMessage{
		Action:        action,
		Payload:       payloadBytes,
		CorrelationID: corrID,
		ReplyTo:       "amq.rabbitmq.reply-to",
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	err = c.channel.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          body,
		},
	)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-respChan:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
