package outbox_processor

import (
	"context"
	"time"

	"github.com/fidesy/sdk/common/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type (
	Service struct {
		topicName string

		storage       Storage
		kafkaProducer KafkaProducer
	}

	Storage interface {
		ListOutboxMessages(ctx context.Context, limit uint64) ([]*Message, error)
		DeleteOutboxMessages(ctx context.Context, ids []int64) error
	}

	KafkaProducer interface {
		ProduceMessage(topic string, messageBytes []byte)
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
	}
}

func (s *Service) Publish(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.publish(ctx)
		}
	}
}

func (s *Service) publish(ctx context.Context) {
	ctx = context.WithValue(ctx, "skip_span", true)

	messages, err := s.storage.ListOutboxMessages(ctx, 100)
	if err != nil {
		logger.Errorf("storage.ListOutboxMessages: %v", err)
		return
	}

	if len(messages) == 0 {
		return
	}

	for _, message := range messages {
		s.kafkaProducer.ProduceMessage(s.topicName, []byte(message.Message))
	}

	err = s.storage.DeleteOutboxMessages(ctx, lo.Map(messages, func(message *Message, _ int) int64 {
		return message.ID
	}))
	if err != nil {
		logger.Errorf("storage.DeleteOutboxMessages: %v", err,
			zap.Int("messageID:", int(messages[0].ID)),
		)
	}
}
