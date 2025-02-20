package outbox_processor

import (
	"context"
	"fmt"
	"time"

	"github.com/fidesy/sdk/common/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type (
	Service struct {
		topicName string

		storage       Storage
		kafkaProducer KafkaProducer

		forcePush chan struct{}
	}

	Storage interface {
		ListOutboxMessages(ctx context.Context, limit uint64) ([]*Message, error)
		DeleteOutboxMessages(ctx context.Context, ids []int64) error
	}

	KafkaProducer interface {
		ProduceMessage(ctx context.Context, messages [][]byte) error
	}
)

func New(
	tableName string,
	topicName string,
	pool *pgxpool.Pool,
	producer KafkaProducer,
) *Service {
	return &Service{
		topicName:     topicName,
		storage:       NewStorage(tableName, pool),
		kafkaProducer: producer,
		forcePush:     make(chan struct{}),
	}
}

func (s *Service) Publish(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.forcePush:
			if err := s.publish(ctx); err != nil {
				logger.Errorf("outbox_processor.forcePush publish: %v", err)
			}
		case <-ticker.C:
			if err := s.publish(ctx); err != nil {
				logger.Errorf("outbox_processor.publish: %v", err)
			}
		}
	}
}

func (s *Service) Push() {
	go func() {
		s.forcePush <- struct{}{}
	}()
}

func (s *Service) publish(ctx context.Context) error {
	ctx = context.WithValue(ctx, "skip_span", true)

	messages, err := s.storage.ListOutboxMessages(ctx, 1000)
	if err != nil {
		return fmt.Errorf("storage.ListOutboxMessages: %w", err)
	}

	if len(messages) == 0 {
		return nil
	}

	kafkaMessages := make([][]byte, 0, len(messages))
	for _, message := range messages {
		kafkaMessages = append(kafkaMessages, []byte(message.Message))
	}

	err = s.kafkaProducer.ProduceMessage(ctx, kafkaMessages)
	if err != nil {
		return fmt.Errorf("kafkaProducer.ProduceMessage: %w", err)
	}

	err = s.storage.DeleteOutboxMessages(ctx, lo.Map(messages, func(message *Message, _ int) int64 {
		return message.ID
	}))
	if err != nil {
		return fmt.Errorf("storage.DeleteOutboxMessages: %w", err)
	}

	return nil
}
