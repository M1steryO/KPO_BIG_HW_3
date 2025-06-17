package services

import (
	"context"
	"log"
	"orders-service/internal/config"
	"strconv"
	"time"

	"orders-service/infrastructure/kafka"
	"orders-service/internal/storage"
)

type Sender struct {
	storage  storage.OrderStorage
	producer kafka.Producer
	kafkaCfg config.KafkaConfig
}

func NewSender(storage storage.OrderStorage, producer kafka.Producer, kafkaCfg config.KafkaConfig) *Sender {
	return &Sender{storage: storage, producer: producer, kafkaCfg: kafkaCfg}
}

func (s *Sender) ProcessOnce() {
	entries, err := s.storage.FetchUnprocessedOutbox()
	if err != nil {
		log.Printf("outbox fetch error: %v", err)
		return
	}
	for _, e := range entries {
		aggregateIdStr := strconv.FormatInt(e.AggregateID, 10)
		if err := s.producer.Publish(s.kafkaCfg.TopicPaymentRequests, aggregateIdStr, e.Payload); err != nil {
			log.Printf("publish outbox id=%d error: %v", e.ID, err)
			continue
		}
		if err := s.storage.MarkOutboxProcessed(e.ID); err != nil {
			log.Printf("mark outbox id=%d processed error: %v", e.ID, err)
		}
	}
}

func (s *Sender) StartProcessEvents(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("stop process events")
				return
			case <-ticker.C:
				s.ProcessOnce()
			}
		}
	}()
}
