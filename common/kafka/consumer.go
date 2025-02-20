package kafka

import (
	"context"
	"errors"
	"github.com/fidesy/sdk/common/logger"
	"github.com/segmentio/kafka-go"
	"math"
	"sync"
	"time"
)

type (
	Consumer struct {
		reader         *kafka.Reader
		messageHandler MessageHandler
		retryConfig    RetryConfig

		locks map[int]*sync.Mutex
	}

	RetryConfig struct {
		// Default math.MaxInt
		MaxRetries int
		// Default 500ms
		RetryDuration time.Duration
	}

	MessageHandler interface {
		ProcessMessage(ctx context.Context, message []byte) error
	}
)

func NewConsumer(
	reader *kafka.Reader,
	messageHandler MessageHandler,
	retryConfig RetryConfig,
	partitionsAmount int,
) *Consumer {
	locks := make(map[int]*sync.Mutex, partitionsAmount)
	for i := 0; i < partitionsAmount; i++ {
		locks[i] = &sync.Mutex{}
	}

	// Apply default config value

	if retryConfig.RetryDuration == time.Duration(0) {
		retryConfig.RetryDuration = 500 * time.Millisecond
	}

	if retryConfig.MaxRetries == 0 {
		retryConfig.MaxRetries = math.MaxInt
	}

	return &Consumer{
		reader:         reader,
		messageHandler: messageHandler,
		retryConfig:    retryConfig,
		locks:          locks,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			logger.Errorf("reader.FetchMessage: %v", err)
			break
		}

		// !!!
		// Every partition in loop

		// Waiting until previous message is done
		c.locks[m.Partition].Lock()

		err = c.messageHandler.ProcessMessage(ctx, m.Value)
		if err != nil {
			logger.Info("error while consuming message...")
			err = c.consumeWithRetries(ctx, m)()
			if err != nil {
				logger.Errorf("consumeWithRetries: %v", err)
				// Must send to DLQ
			}
		}

		err = c.reader.CommitMessages(ctx, m)
		if err != nil {
			logger.Errorf("reader.CommitMessages: %v", err)
		}

		c.locks[m.Partition].Unlock()
	}

	return nil
}

func (c *Consumer) consumeWithRetries(ctx context.Context, message kafka.Message) func() error {
	retries := 0

	return func() error {
		for retries < c.retryConfig.MaxRetries {
			retries++

			err := c.messageHandler.ProcessMessage(ctx, message.Value)
			if err != nil {
				logger.Errorf("messageHandler.ProcessMessage: %v", err)
				time.Sleep(c.retryConfig.RetryDuration)
				continue
			}

			return nil
		}

		return errors.New("error while consuming message, all attempts exceeded")
	}
}
