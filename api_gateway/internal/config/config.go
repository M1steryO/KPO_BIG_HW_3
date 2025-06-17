package config

import (
	"errors"
	"os"
)

type Config struct {
	GatewayAddr       string
	OrderServiceURL   string
	PaymentServiceURL string
}

func MustLoad() (*Config, error) {
	cfg := &Config{
		GatewayAddr:       os.Getenv("HTTP_ADDR"),
		OrderServiceURL:   os.Getenv("ORDER_SERVICE_URL"),
		PaymentServiceURL: os.Getenv("PAYMENT_SERVICE_URL"),
	}
	if cfg.OrderServiceURL == "" || cfg.PaymentServiceURL == "" {
		return nil, errors.New("ORDER_SERVICE_URL and PAYMENT_SERVICE_URL are required")
	}
	return cfg, nil
}
