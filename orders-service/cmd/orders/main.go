package main

import (
	"context"
	"log"
	"orders-service/infrastructure/kafka"
	"orders-service/internal/api/http"
	"orders-service/internal/config"
	"orders-service/internal/services"
	"orders-service/internal/storage"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("kafka producer init failed: %v", err)
	}

	orderStore, err := storage.NewPostgresOrderStorage(cfg.DB)
	if err != nil {
		log.Fatalf("db init failed: %v", err)

	}

	sender := services.NewSender(orderStore, producer, cfg.Kafka)
	sender.StartProcessEvents(context.Background())

	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicPaymentProcessed,
		"order-service-inbox",
	)
	if err != nil {
		log.Fatalf("kafka consumer init: %v", err)
	}
	processor := services.NewPaymentResultProcessor(
		orderStore,
		consumer,
		cfg.Kafka.TopicPaymentProcessed,
	)
	processor.Start(context.Background())

	router := http.NewOrderRouter(orderStore)
	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
	log.Printf("server shutdown")
}
