package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

type BatchHandler[T any] struct {
	l     logger.Logger
	bizFn func(msgs []*sarama.ConsumerMessage, events []T) error
}

func NewBatchHandler[T any](
	l logger.Logger,
	bizFn func(msgs []*sarama.ConsumerMessage, events []T) error,
) *BatchHandler[T] {
	return &BatchHandler[T]{
		l:     l,
		bizFn: bizFn,
	}
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *BatchHandler[T]) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *BatchHandler[T]) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		events := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		done := false
		for range batchSize {
			select {
			case <-ctx.Done():
				// timeout
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				var event T
				err := json.Unmarshal(msg.Value, &event)
				if err != nil {
					h.l.Error(
						"failed to unmarshal",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err),
					)
					continue
				}
				batch = append(batch, msg)
				events = append(events, event)
			}
			if done {
				break
			}
		}
		cancel()

		// batch process
		err := h.bizFn(batch, events)
		if err != nil {
			// NOTE: can introduce a retry
			h.l.Error(
				"failed to handle msgs batch",
				logger.Error(err),
			)
			continue
		}

		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}

// Setup implements sarama.ConsumerGroupHandler.
func (h *BatchHandler[T]) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
