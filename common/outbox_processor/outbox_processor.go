package outbox_processor

import (
	"context"
	"fmt"
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

		forcePush chan struct{}
	}

	Storage interface {
		ListOutboxMessages(ctx context.Context, limit int64) ([]*Message, error)
		UpdateOutboxMessagesSentAt(ctx context.Context, ids []int64) error
	}

	KafkaProducer interface {
		ProduceMessage(ctx context.Context, messages [][]byte) error
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
		forcePush:     make(chan struct{}),
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

	messages, err := s.storage.ListOutboxMessages(ctx, s.listOutboxLimit)
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

	messagesIDs := lo.Map(messages, func(message *Message, _ int) int64 {
		return message.ID
	})

	err = s.storage.UpdateOutboxMessagesSentAt(ctx, messagesIDs)
	if err != nil {
		logger.Errorf("storage.UpdateOutboxMessagesSentAt: %v", err,
			zap.Int("messageID:", int(messages[0].ID)),
		)
	}

	return nil
}
