// Package worker provides background task and message consumption workers.
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
)

// StartNotificationWorker starts consuming notification messages from the message broker.
func StartNotificationWorker(ctx context.Context, broker port.MessageBroker) {
	timeout := 10 * time.Second

	err := broker.Consume(ctx, "notification.queue", func(ctx context.Context, body []byte) error {
		data := map[string]any{}
		if err := sonic.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unmarshal failed: %w", err)
		}

		log.Printf("Processing event: %v", data["event"])

		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}

		log.Println("Success processing notification")
		return nil
	}, timeout)
	if err != nil {
		log.Printf("notification worker failed to start: %v", err)
	}
}
