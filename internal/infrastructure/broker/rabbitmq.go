// Package broker provides RabbitMQ implementation for message publishing and consumption.
package broker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQImpl struct {
	conn           *amqp.Connection
	publishChannel *amqp.Channel
	consumeChannel *amqp.Channel
	mu             sync.RWMutex
	done           chan struct{}
}

// NewRabbitMQ creates a new RabbitMQ message broker instance and declares exchanges, queues, and DLX topologies.
func NewRabbitMQ(cfg *config.Config) (port.MessageBroker, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}

	publishChannel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open publish channel: %w", err)
	}

	if err := publishChannel.Confirm(false); err != nil {
		_ = publishChannel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	consumeChannel, err := conn.Channel()
	if err != nil {
		_ = publishChannel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open consume channel: %w", err)
	}

	// Configure QoS (Prefetch Count) so that RabbitMQ does not overwhelm the consumer all at once.
	if err := consumeChannel.Qos(1, 0, false); err != nil {
		_ = consumeChannel.Close()
		_ = publishChannel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	if err := declareTopology(consumeChannel); err != nil {
		_ = consumeChannel.Close()
		_ = publishChannel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare topology: %w", err)
	}

	log.Println("RabbitMQ successfully connected. DLX topology ready.")

	return &rabbitMQImpl{
		conn:           conn,
		publishChannel: publishChannel,
		consumeChannel: consumeChannel,
		done:           make(chan struct{}),
	}, nil
}

func declareTopology(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare("dlx.exchange", "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare dlx exchange: %w", err)
	}

	if _, err := ch.QueueDeclare("dlx.queue", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare dlx queue: %w", err)
	}

	if err := ch.QueueBind("dlx.queue", "dlx.routingKey", "dlx.exchange", false, nil); err != nil {
		return fmt.Errorf("failed to bind dlx queue: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "dlx.exchange",
		"x-dead-letter-routing-key": "dlx.routingKey",
	}

	if _, err := ch.QueueDeclare("notification.reset_password", true, false, false, false, args); err != nil {
		return fmt.Errorf("failed to declare notification queue: %w", err)
	}

	return nil
}

// Publish sends a message to the specified exchange with a given routing key.
func (r *rabbitMQImpl) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	r.mu.RLock()
	ch := r.publishChannel
	r.mu.RUnlock()

	conf, err := ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if !conf.Wait() {
		return fmt.Errorf("message nacked by broker")
	}

	return nil
}

// Consume starts listening to a queue and handles messages asynchronously with timeout support.
func (r *rabbitMQImpl) Consume(
	ctx context.Context,
	queue string,
	handler func(ctx context.Context, body []byte) error,
	timeout time.Duration,
) error {
	r.mu.RLock()
	ch := r.consumeChannel
	r.mu.RUnlock()

	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[RABBITMQ] Consumer stopped.")
				return
			case msg, ok := <-msgs:
				if !ok {
					return // channel closed
				}

				msgCtx, cancel := context.WithTimeout(ctx, timeout)

				errCh := make(chan error, 1)
				go func(m amqp.Delivery) {
					errCh <- handler(msgCtx, m.Body)
				}(msg)

				select {
				case <-msgCtx.Done():
					log.Printf("[RABBITMQ] Timeout (%v). Nack → DLX.", timeout)
					if err := msg.Nack(false, false); err != nil {
						log.Printf("[RABBITMQ] Failed to Nack message on timeout: %v\n", err)
					}
				case err := <-errCh:
					if err != nil {
						log.Printf("[RABBITMQ] Error: %v. Nack → DLX.", err)
						if err := msg.Nack(false, false); err != nil {
							log.Printf("[RABBITMQ] Failed to Nack message on error: %v\n", err)
						}
					} else {
						if err := msg.Ack(false); err != nil {
							log.Printf("[RABBITMQ] Failed to Ack message: %v\n", err)
						}
					}
				}
				cancel()
			}
		}
	}()

	return nil
}

// Close closes the RabbitMQ channel and connection gracefully.
func (r *rabbitMQImpl) Close() error {
	close(r.done)
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.consumeChannel != nil {
		if err := r.consumeChannel.Close(); err != nil {
			return fmt.Errorf("failed to close consume channel: %w", err)
		}
	}

	if r.publishChannel != nil {
		if err := r.publishChannel.Close(); err != nil {
			return fmt.Errorf("failed to close publish channel: %w", err)
		}
	}

	return r.conn.Close()
}
