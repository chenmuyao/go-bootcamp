package ioc

import (
	"github.com/IBM/sarama"
	"github.com/chenmuyao/go-bootcamp/interactive/config"
	intrEvents "github.com/chenmuyao/go-bootcamp/interactive/events"
	"github.com/chenmuyao/go-bootcamp/internal/events"
)

func InitSaramaClient() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(config.Cfg.Sarama.Addr, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(c1 *intrEvents.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{c1}
}
