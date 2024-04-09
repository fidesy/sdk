package kafka

import (
	"context"
	"fmt"
	"github.com/fidesy/sdk/common/logger"
)

type MessageHandler interface {
	Consume(ctx context.Context, msg []byte) error
}

func RegisterConsumer(
	ctx context.Context,
	handler MessageHandler,
	kafkaBrokers []string,
	topic string,
) error {
	kafkaConsumer, err := NewConsumer(ctx, kafkaBrokers, topic)
	if err != nil {
		return fmt.Errorf("kafka.NewConsumer: %w", err)
	}

	go consume(ctx, handler, kafkaConsumer)

	return nil
}

func consume(ctx context.Context, handler MessageHandler, consumer *Consumer) {
	defer func() {
		err := consumer.Close()
		if err != nil {
			logger.Errorf("kafkaConsumer.Close: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-consumer.Consume():
			err := handler.Consume(ctx, msg.Value)
			if err != nil {
				logger.Errorf("handler.Consume: %v", err)
			}
		}
	}
}
