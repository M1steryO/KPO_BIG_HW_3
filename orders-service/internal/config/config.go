package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type DBConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Config struct {
	DB               DBConfig
	Kafka            KafkaConfig
	HTTPAddr         string
	LogLevel         string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
}

type KafkaConfig struct {
	Brokers               string
	TopicPaymentRequests  string
	TopicPaymentProcessed string
}

func MustLoad() (*Config, error) {
	var err error
	cfg := &Config{}

	cfg.DB.URL = os.Getenv("DATABASE_URL")
	if cfg.DB.URL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if s := os.Getenv("DB_MAX_OPEN_CONNS"); s != "" {
		cfg.DB.MaxOpenConns, err = strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
		}
	} else {
		cfg.DB.MaxOpenConns = 25
	}

	if s := os.Getenv("DB_MAX_IDLE_CONNS"); s != "" {
		cfg.DB.MaxIdleConns, err = strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
		}
	} else {
		cfg.DB.MaxIdleConns = 25
	}

	if s := os.Getenv("DB_CONN_MAX_LIFETIME"); s != "" {
		cfg.DB.ConnMaxLifetime, err = time.ParseDuration(s)
		if err != nil {
			return nil, fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %w", err)
		}
	} else {
		cfg.DB.ConnMaxLifetime = 5 * time.Minute
	}

	cfg.Kafka.Brokers = os.Getenv("KAFKA_BROKERS")
	if cfg.Kafka.Brokers == "" {
		return nil, errors.New("KAFKA_BROKERS is required")
	}

	cfg.Kafka.TopicPaymentRequests = os.Getenv("KAFKA_TOPIC_PAYMENT_REQUESTS")
	if cfg.Kafka.TopicPaymentRequests == "" {
		return nil, errors.New("KAFKA_TOPIC_PAYMENT_PROCESSED is required")
	}

	cfg.Kafka.TopicPaymentProcessed = os.Getenv("KAFKA_TOPIC_PAYMENT_PROCESSED")
	if cfg.Kafka.TopicPaymentProcessed == "" {
		return nil, errors.New("KAFKA_TOPIC_PAYMENT_PROCESSED is required")
	}

	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		cfg.HTTPAddr = addr
	} else {
		cfg.HTTPAddr = ":8080"
	}

	cfg.LogLevel = os.Getenv("LOG_LEVEL")

	if s := os.Getenv("HTTP_READ_TIMEOUT"); s != "" {
		cfg.HTTPReadTimeout, err = time.ParseDuration(s)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
		}
	} else {
		cfg.HTTPReadTimeout = 10 * time.Second
	}
	if s := os.Getenv("HTTP_WRITE_TIMEOUT"); s != "" {
		cfg.HTTPWriteTimeout, err = time.ParseDuration(s)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
		}
	} else {
		cfg.HTTPWriteTimeout = 10 * time.Second
	}

	return cfg, nil
}
