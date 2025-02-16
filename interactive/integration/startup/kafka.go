package startup

import (
	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/interactive/events"
)

func InitSaramaClient() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(c1 *events.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{c1}
}
