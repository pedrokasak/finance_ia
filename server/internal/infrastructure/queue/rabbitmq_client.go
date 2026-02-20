package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Queue names
const (
	QueueAIInsight          = "ai.insight.generate"
	QueueEmailSend          = "email.send"
	QueueStripeWebhook      = "stripe.webhook"
	QueueBehavioralAnalysis = "behavioral.analysis"
)

// Message holds a queue message payload
type Message struct {
	Queue   string          `json:"queue"`
	Payload json.RawMessage `json:"payload"`
}

// Client manages RabbitMQ connection and channel
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient connects to RabbitMQ using RABBITMQ_URL env var
func NewClient() (*Client, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: connection failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: channel creation failed: %w", err)
	}

	client := &Client{conn: conn, channel: ch}

	// Declare all queues
	queues := []string{
		QueueAIInsight,
		QueueEmailSend,
		QueueStripeWebhook,
		QueueBehavioralAnalysis,
	}
	for _, q := range queues {
		if err := client.declareQueue(q); err != nil {
			return nil, fmt.Errorf("rabbitmq: declare queue %s: %w", q, err)
		}
	}

	return client, nil
}

func (c *Client) declareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	return err
}

// Publish sends a message to a queue
func (c *Client) Publish(ctx context.Context, queue string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rabbitmq: marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.channel.PublishWithContext(ctx,
		"",    // exchange (default)
		queue, // routing key = queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // survives broker restart
			Body:         data,
		},
	)
}

// Consume starts consuming messages from a queue
// handler receives the raw JSON bytes and returns error (message is nacked on error)
func (c *Client) Consume(queue string, handler func([]byte) error) error {
	msgs, err := c.channel.Consume(
		queue,
		"",    // consumer tag
		false, // auto-ack (manual for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq: consume %s: %w", queue, err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("[queue] %s: handler error: %v — nacking", queue, err)
				_ = d.Nack(false, true) // requeue
			} else {
				_ = d.Ack(false)
			}
		}
	}()

	return nil
}

// Close gracefully closes the connection
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
