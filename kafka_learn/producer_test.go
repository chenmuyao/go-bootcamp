package kafkalearn

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

var addr = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	// NOTE: must have
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(addr, cfg)
	// cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	// cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	// // by key
	// cfg.Producer.Partitioner = sarama.NewHashPartitioner
	// // use the partition in SendMessage
	// cfg.Producer.Partitioner = sarama.NewManualPartitioner
	// // consistent hash
	// cfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	// // custom hash
	// cfg.Producer.Partitioner = sarama.NewCustomPartitioner()
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("this is a message"),

		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("heanders are transferred"),
				Value: []byte("from prudcer to consumer"),
			},
		},
		Metadata: "Metadata is not transferred from producer to consumer",
	})
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addr, cfg)
	assert.NoError(t, err)
	msgs := producer.Input()

	msgs <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("this is a message"),

		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("heanders are transferred"),
				Value: []byte("from prudcer to consumer"),
			},
		},
		Metadata: "Metadata is not transferred from producer to consumer",
	}

	select {
	case msg := <-producer.Successes():
		t.Log("send sucess", string(msg.Value.(sarama.StringEncoder)))
	case err := <-producer.Errors():
		t.Log("send fail", err.Error())
	}
}
