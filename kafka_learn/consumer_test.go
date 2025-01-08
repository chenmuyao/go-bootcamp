package kafkalearn

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "test", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	err = consumer.Consume(ctx, []string{"test_topic"}, &ConsumerHandler{})
	assert.NoError(t, err)
}

type ConsumerHandler struct{}

// Cleanup implements sarama.ConsumerGroupHandler.
func (c *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("cleanup")
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (c *ConsumerHandler) ConsumeClaimSingleMsg(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

// Batch
func (c *ConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		done := false
		var eg errgroup.Group
		for i := range batchSize {
			select {
			case <-ctx.Done():
				// timeout
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch[i] = msg
				eg.Go(func() error {
					log.Println(msg)
					return nil
				})
			}
			if done {
				break
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			log.Fatal(err)
			continue
		}

		// batch process
		// log.Println(batch)

		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}

// Setup implements sarama.ConsumerGroupHandler.
func (c *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("setup")
	var offset int64 = 0
	// sarama.OffsetOldest
	// sarama.OffsetNewest
	partitions := session.Claims()["test_topic"]
	for _, p := range partitions {
		session.ResetOffset("test_topic", p, offset, "")
	}
	return nil
}

var _ sarama.ConsumerGroupHandler = &ConsumerHandler{}
