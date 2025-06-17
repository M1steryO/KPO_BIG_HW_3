package services

import (
	"context"
	"log"
	"payments-service/internal/config"
	"time"

	"payments-service/infrastructure/kafka"
	"payments-service/internal/storage"
)

type OutboxSender struct {
	store    storage.AccountStorage
	producer kafka.Producer
	kafka    config.KafkaConfig
}

func NewOutboxSender(store storage.AccountStorage, producer kafka.Producer, kafka config.KafkaConfig) *OutboxSender {
	return &OutboxSender{store: store, producer: producer, kafka: kafka}
}

func (o *OutboxSender) StartProcessEvents(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("outbox sender stopped")
				return
			case <-ticker.C:
				entries, err := o.store.FetchUnprocessedOutbox()
				if err != nil {
					log.Printf("fetch outbox err: %v", err)
					continue
				}
				for _, e := range entries {
					if err := o.producer.Publish(o.kafka.TopicPaymentProcessed, e.EventType, e.Message, e.Payload); err != nil {
						log.Printf("publish payment_processed id=%d err: %v", e.ID, err)
						continue
					}
					if err := o.store.MarkOutboxProcessed(e.ID); err != nil {
						log.Printf("mark processed id=%d err: %v", e.ID, err)
					}
				}
			}
		}
	}()
}
