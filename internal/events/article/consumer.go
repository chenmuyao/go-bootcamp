package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/chenmuyao/go-bootcamp/pkg/saramax"
	"github.com/prometheus/client_golang/prometheus"
)

const consumeTimeout = time.Second

type Consumer interface {
	Consume(msg *sarama.ConsumerMessage, event ReadEvent) error
}

type InteractiveReadEventConsumer struct {
	l      logger.Logger
	repo   repository.InteractiveRepository
	client sarama.Client
}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{TopicReadEvent},
			saramax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume),
		)
		if er != nil {
			i.l.Error("quit consuming", logger.Error(er))
		}
	}()
	return nil
}

// StartV1 consume one message a time
func (i *InteractiveReadEventConsumer) StartV1() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	promOpts := prometheus.SummaryOpts{
		Namespace: "my_company",
		Subsystem: "wetravel",
		Name:      "kafka_consumer",
		Help:      "Kafka consumer metrics",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
		ConstLabels: prometheus.Labels{
			"instance_id": "instance",
		},
	}
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{TopicReadEvent},
			saramax.NewHandler[ReadEvent](i.l, promOpts, i.Consume),
		)
		if er != nil {
			i.l.Error("quit consuming", logger.Error(er))
		}
	}()
	return nil
}

// Consume implements Consumer.
func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), consumeTimeout)
	defer cancel()

	i.l.Debug("Consume")

	return i.repo.IncrReadCnt(ctx, "article", event.Aid)
}

func (i *InteractiveReadEventConsumer) BatchConsume(
	msgs []*sarama.ConsumerMessage,
	events []ReadEvent,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), consumeTimeout)
	defer cancel()

	bizs := make([]string, len(events))
	bizIDs := make([]int64, len(events))

	for i, ev := range events {
		bizs[i] = "article"
		bizIDs[i] = ev.Aid
	}

	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIDs)
}

func NewInteractiveReadEventConsumer(
	l logger.Logger,
	repo repository.InteractiveRepository,
	client sarama.Client,
) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		l:      l,
		repo:   repo,
		client: client,
	}
}
