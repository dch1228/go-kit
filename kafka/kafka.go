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
}

type subscriber struct {
	cg      sarama.ConsumerGroup
	kConfig *sarama.Config
	topics  []string
	groupID string
	handler *consumerGroupHandler
}

type Handler func(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage)

type consumerGroupHandler struct {
	handler Handler
}

func New(cfg Config) *Kafka {
	return &Kafka{
		cfg:         cfg,
		log:         log.L().Named("[kafka]"),
		subscribers: make([]*subscriber, 0, 10),
	}
}

func (c *Kafka) Subscribe(cfg ConsumerConfig, handler Handler) {
	kConfig := sarama.NewConfig()

	c.subscribers = append(c.subscribers, &subscriber{
		kConfig: kConfig,
		topics:  []string{cfg.Topic},
		groupID: cfg.GroupID,
		handler: &consumerGroupHandler{
			handler: handler,
		},
	})
}

func (c *Kafka) Start(_ context.Context) error {
	eg := errgroup.Group{}
	for _, subscriber := range c.subscribers {
		subscriber := subscriber
		eg.Go(func() error {
			cg, err := sarama.NewConsumerGroup(
				c.cfg.BootstrapServers,
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

func (c *Kafka) Stop(_ context.Context) error {
	for _, subscriber := range c.subscribers {
		c.log.Info("closing consumers", log.String("topic", subscriber.topics[0]))
		if err := subscriber.cg.Close(); err != nil {
			c.log.Error("close consume error", err)
		}
	}
	return nil
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(session, msg)
	}
	return nil
}
