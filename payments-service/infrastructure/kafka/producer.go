package kafka

import (
	_ "fmt"
	"github.com/IBM/sarama"
	"strings"
	"time"
)

type Producer interface {
	Publish(topic string, eventType, message string, payload []byte) error
}

type producer struct {
	client sarama.SyncProducer
}

func NewProducer(brokers string) (Producer, error) {
	addrs := strings.Split(brokers, ",")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &producer{client: client}, nil
}

func (p *producer) Publish(topic string, eventType, message string, payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(eventType)},
			{Key: []byte("message"), Value: []byte(message)},
		},
		Timestamp: time.Now(),
	}
	_, _, err := p.client.SendMessage(msg)
	return err
}
