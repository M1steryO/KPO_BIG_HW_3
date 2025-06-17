package services

import (
	"context"
	"log"

	"payments-service/infrastructure/kafka"
	"payments-service/internal/storage"
)

type InboxProcessor struct {
	store    storage.AccountStorage
	consumer kafka.Consumer
}

func NewInboxProcessor(store storage.AccountStorage, consumer kafka.Consumer) *InboxProcessor {
	return &InboxProcessor{store: store, consumer: consumer}
}

func (p *InboxProcessor) StartProcess(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("inbox processor stopped")
				return
			case msg := <-p.consumer.Messages():
				key := string(msg.Key)
				inboxId, err := p.store.ProcessPaymentRequest(key, msg.Value)
				if err != nil {
					log.Printf("process payment request key=%s error: %v", key, err)
					continue
				}
				if err = p.store.MarkInboxProcessed(inboxId); err != nil {
					log.Printf("failed to mark inbox %d processed: %v", inboxId, err)
				}
				p.consumer.MarkProcessed(msg)
			}
		}
	}()
}
