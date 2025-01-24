package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/chenmuyao/go-bootcamp/pkg/saramax"
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
			saramax.NewHandler[ReadEvent](i.l, i.Consume),
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
