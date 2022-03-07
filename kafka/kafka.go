package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"golang.org/x/sync/errgroup"

	"github.com/dch1228/go-kit/log"
)

type Kafka struct {
	cfg        Config
	log        *log.Logger
	consumers  []*consumer
	middleware []MiddlewareFunc
}

func New(cfg Config) *Kafka {
	return &Kafka{
		cfg:       cfg,
		log:       log.L().Named("[kafka]"),
		consumers: make([]*consumer, 0, 10),
	}
}

func (k *Kafka) Use(middleware ...MiddlewareFunc) {
	k.middleware = append(k.middleware, middleware...)
}

func (k *Kafka) Subscribe(cfg ConsumerConfig, handler HandlerFunc) {
	kConfig := sarama.NewConfig()

	k.consumers = append(k.consumers, &consumer{
		kConfig: kConfig,
		topics:  []string{cfg.Topic},
		groupID: cfg.GroupID,
		handler: &consumerGroupHandler{
			k:       k,
			handler: handler,
		},
	})
}

func (k *Kafka) Start(_ context.Context) error {
	eg := errgroup.Group{}
	for _, consumer := range k.consumers {
		consumer := consumer
		eg.Go(func() error {
			cg, err := sarama.NewConsumerGroup(
				k.cfg.BootstrapServers,
				consumer.groupID,
				consumer.kConfig,
			)
			if err != nil {
				return err
			}
			ctx := context.Background()
			consumer.cg = cg
			for {
				if err := cg.Consume(ctx, consumer.topics, consumer.handler); err != nil {
					return err
				}
			}
		})
	}

	return eg.Wait()
}

func (k *Kafka) Stop(_ context.Context) error {
	for _, consumer := range k.consumers {
		k.log.Info("closing consumer", log.String("topic", consumer.topics[0]))
		if err := consumer.cg.Close(); err != nil {
			k.log.Error("close consume error", err, log.String("topic", consumer.topics[0]))
		}
	}
	return nil
}
