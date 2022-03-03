package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"golang.org/x/sync/errgroup"

	"github.com/dch1228/go-kit/log"
)

type Kafka struct {
	cfg         Config
	log         *log.Logger
	subscribers []*subscriber
	middleware  []MiddlewareFunc
}

func New(cfg Config) *Kafka {
	return &Kafka{
		cfg:         cfg,
		log:         log.L().Named("[kafka]"),
		subscribers: make([]*subscriber, 0, 10),
	}
}

func (k *Kafka) Use(middleware ...MiddlewareFunc) {
	k.middleware = append(k.middleware, middleware...)
}

func (k *Kafka) Subscribe(cfg ConsumerConfig, handler HandlerFunc) {
	kConfig := sarama.NewConfig()

	k.subscribers = append(k.subscribers, &subscriber{
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
	for _, subscriber := range k.subscribers {
		subscriber := subscriber
		eg.Go(func() error {
			cg, err := sarama.NewConsumerGroup(
				k.cfg.BootstrapServers,
				subscriber.groupID,
				subscriber.kConfig,
			)
			if err != nil {
				return err
			}
			ctx := context.Background()
			subscriber.cg = cg
			for {
				if err := cg.Consume(ctx, subscriber.topics, subscriber.handler); err != nil {
					return err
				}
			}
		})
	}

	return eg.Wait()
}

func (k *Kafka) Stop(_ context.Context) error {
	for _, subscriber := range k.subscribers {
		k.log.Info("closing consumers", log.String("topic", subscriber.topics[0]))
		if err := subscriber.cg.Close(); err != nil {
			k.log.Error("close consume error", err, log.String("topic", subscriber.topics[0]))
		}
	}
	return nil
}
