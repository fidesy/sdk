package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(
	writer *kafka.Writer,
) *Producer {
	return &Producer{
		writer: writer,
	}
}

func (p *Producer) ProduceMessage(ctx context.Context, messages [][]byte) error {
	kafkaMessages := make([]kafka.Message, 0, len(messages))
	for _, msg := range messages {
		kafkaMessages = append(kafkaMessages, kafka.Message{
			Value: msg,
		})
	}

	err := p.writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		return fmt.Errorf("writer.WriteMessages: %w", err)
	}

	return nil
}
