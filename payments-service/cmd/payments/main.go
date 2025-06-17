package main

import (
	"context"
	"log"
	"payments-service/infrastructure/kafka"
	"payments-service/internal/api/http"
	"payments-service/internal/config"
	"payments-service/internal/services"
	"payments-service/internal/storage"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	store, err := storage.NewPostgresAccountStorage(cfg.DB)
	if err != nil {
		log.Fatalf("storage init failed: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("kafka producer init: %v", err)
	}
	outboxSender := services.NewOutboxSender(store, producer, cfg.Kafka)
	go outboxSender.StartProcessEvents(context.Background())

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers,
		cfg.Kafka.TopicPaymentRequests, "payments-service-inbox")
	if err != nil {
		log.Fatalf("kafka consumer init: %v", err)
	}
	inboxProcessor := services.NewInboxProcessor(store, consumer)
	inboxProcessor.StartProcess(context.Background())

	r := http.NewAccountRouter(store)
	log.Printf("HTTP listening on %s", cfg.HTTPAddr)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
	log.Printf("server shutdown!")
}
