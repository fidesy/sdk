package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(ctx context.Context, kafkaBrokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(kafkaBrokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewAsyncProducer: %w", err)
	}

	kafkaProducer := &Producer{producer: producer}

	return kafkaProducer, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) ProduceMessage(topic string, messageBytes []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	p.producer.Input() <- msg
}
