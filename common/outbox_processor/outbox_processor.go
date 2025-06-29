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
		fetchDuration   time.Duration
		listOutboxLimit int64

		topicName string

		storage       Storage
		kafkaProducer KafkaProducer
	}

	Storage interface {
		ListOutboxMessages(ctx context.Context, limit int64) ([]*Message, error)
		UpdateOutboxMessagesSentAt(ctx context.Context, ids []int64) error
	}

	KafkaProducer interface {
		ProduceMessage(topic string, messageBytes []byte)
	}
)

type Option func(s *Service)

func WithFetchDuration(duration time.Duration) Option {
	return func(s *Service) {
		s.fetchDuration = duration
	}
}

func WithListOutboxLimit(limit int64) Option {
	return func(s *Service) {
		s.listOutboxLimit = limit
	}
}

func New(
	tableName string,
	topicName string,
	pool *pgxpool.Pool,
	producer KafkaProducer,
	options ...Option,
) *Service {
	s := &Service{
		fetchDuration:   500 * time.Millisecond,
		listOutboxLimit: 100,

		topicName:     topicName,
		storage:       NewStorage(tableName, pool),
		kafkaProducer: producer,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Service) Publish(ctx context.Context) {
	ticker := time.NewTicker(s.fetchDuration)
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

	messages, err := s.storage.ListOutboxMessages(ctx, s.listOutboxLimit)
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

	messagesIDs := lo.Map(messages, func(message *Message, _ int) int64 {
		return message.ID
	})

	err = s.storage.UpdateOutboxMessagesSentAt(ctx, messagesIDs)
	if err != nil {
		logger.Errorf("storage.UpdateOutboxMessagesSentAt: %v", err,
			zap.Int("messageID:", int(messages[0].ID)),
		)
	}
}
