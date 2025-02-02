package saramax

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler[T any] struct {
	l      logger.Logger
	bizFn  func(msg *sarama.ConsumerMessage, event T) error
	vector *prometheus.SummaryVec
}

func NewHandler[T any](
	l logger.Logger,
	vector *prometheus.SummaryVec,
	bizFn func(msg *sarama.ConsumerMessage, event T) error,
) *Handler[T] {
	return &Handler[T]{
		l:      l,
		vector: vector,
		bizFn:  bizFn,
	}
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *Handler[T]) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *Handler[T]) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	msgs := claim.Messages()
	for msg := range msgs {
		h.handle(msg, session)
	}
	return nil
}

func (h *Handler[T]) handle(
	msg *sarama.ConsumerMessage,
	session sarama.ConsumerGroupSession,
) {
	var t T
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		h.vector.WithLabelValues(msg.Topic, strconv.Itoa(int(msg.Partition))).
			Observe(float64(duration))
	}()
	err := json.Unmarshal(msg.Value, &t)
	if err != nil {
		// NOTE: can introduce a retry
		h.l.Error(
			"failed to unmarshal",
			logger.String("topic", msg.Topic),
			logger.Int32("partition", msg.Partition),
			logger.Int64("offset", msg.Offset),
			logger.Error(err),
		)
	}
	// run biz
	err = h.bizFn(msg, t)
	if err != nil {
		// NOTE: can introduce a retry
		h.l.Error(
			"failed to handle msg",
			logger.String("topic", msg.Topic),
			logger.Int32("partition", msg.Partition),
			logger.Int64("offset", msg.Offset),
			logger.Error(err),
		)
	}
	session.MarkMessage(msg, "")
}

// Setup implements sarama.ConsumerGroupHandler.
func (h *Handler[T]) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
