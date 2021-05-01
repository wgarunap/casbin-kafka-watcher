package watcher

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/tryfix/kstream/consumer"
	"github.com/tryfix/log"
)

func newConsumer(brokers  []string, logger log.Logger) consumer.Consumer {
	config := consumer.NewConsumerConfig()
	config.BootstrapServers = brokers
	config.GroupId = "casbin-consumer-group"
	config.Logger = logger
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.DefaultVersion
	consumerOb, err := consumer.NewConsumer(config)
	if err != nil {
		log.Fatal(err)
	}

	return consumerOb

}

type rebalancedHandler struct{}

func (r rebalancedHandler) OnPartitionRevoked(ctx context.Context, revoked []consumer.TopicPartition) error {
	return nil
}

func (r rebalancedHandler) OnPartitionAssigned(ctx context.Context, assigned []consumer.TopicPartition) error {
	return nil
}
