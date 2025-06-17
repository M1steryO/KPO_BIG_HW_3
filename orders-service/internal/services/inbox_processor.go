package services

import (
	"context"
	"encoding/json"
	"log"
	"orders-service/infrastructure/kafka"
	"orders-service/internal/models"
	"orders-service/internal/storage"
)

type PaymentResultProcessor struct {
	orders   storage.OrderStorage
	consumer kafka.Consumer
	topic    string
}

func NewPaymentResultProcessor(
	orders storage.OrderStorage,
	consumer kafka.Consumer,
	topic string,
) *PaymentResultProcessor {
	return &PaymentResultProcessor{
		orders: orders, consumer: consumer,
		topic: topic,
	}
}

const (
	paymentStatusFailed  = "failed"
	paymentStatusSuccess = "success"
)

func (p *PaymentResultProcessor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("payment result processor stopped")
				return
			case msg := <-p.consumer.Messages():
				var order models.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("invalid event payload: %v", err)
					p.consumer.MarkProcessed(msg)
					continue
				}

				var newStatus models.Status
				var eventType string
				var errorMessage string
				for _, h := range msg.Headers {
					if string(h.Key) == "event_type" {
						eventType = string(h.Value)
					}
					if string(h.Key) == "message" {
						errorMessage = string(h.Value)
					}
				}
				switch eventType {
				case paymentStatusFailed:
					newStatus = models.Cancelled
					break
				case paymentStatusSuccess:
					newStatus = models.Finished
					break
				default:
					newStatus = models.Cancelled
				}

				if err := p.orders.UpdateStatus(order.ID, newStatus, errorMessage); err != nil {
					log.Printf("failed to update order %d status to %s: %v", order.ID, newStatus, err)
				}
				p.consumer.MarkProcessed(msg)
			}
		}
	}()
}
