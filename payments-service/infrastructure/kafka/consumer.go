package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"strings"
)

type Consumer interface {
	Messages() <-chan *sarama.ConsumerMessage
	MarkProcessed(msg *sarama.ConsumerMessage)
}

type consumer struct {
	pc       sarama.PartitionConsumer
	messages chan *sarama.ConsumerMessage
	pom      sarama.PartitionOffsetManager
}

func NewConsumer(brokers, topic, group string) (Consumer, error) {
	addrs := strings.Split(brokers, ",")

	client, err := sarama.NewClient(addrs, nil)
	if err != nil {
		return nil, fmt.Errorf("kafka NewClient failed: %w", err)
	}

	om, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return nil, fmt.Errorf("offset manager init failed: %w", err)
	}

	pom, err := om.ManagePartition(topic, 0)
	if err != nil {
		return nil, fmt.Errorf("manage partition failed: %w", err)
	}

	consumerClient, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("consumer from client failed: %w", err)
	}

	pc, err := consumerClient.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("consume partition failed: %w", err)
	}

	ch := make(chan *sarama.ConsumerMessage)
	go func() {
		defer close(ch)
		for msg := range pc.Messages() {
			ch <- msg
		}
	}()

	return &consumer{pc: pc, messages: ch, pom: pom}, nil
}

func (c *consumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}

func (c *consumer) MarkProcessed(msg *sarama.ConsumerMessage) {
	offset := msg.Offset + 1
	c.pom.MarkOffset(offset, "")
}
