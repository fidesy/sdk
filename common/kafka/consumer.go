package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.PartitionConsumer
}

func NewConsumer(_ context.Context, brokers []string, topicName string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumer: %w", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("consumer.ConsumePartition: %w", err)
	}

	kafkaConsumer := &Consumer{consumer: partitionConsumer}

	return kafkaConsumer, nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) Consume() <-chan *sarama.ConsumerMessage {
	return c.consumer.Messages()
}
