package kafka

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"

	"github.com/dch1228/go-kit/log"
)

func mockHandler(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	log.Info("mock handler",
		log.String("topic", msg.Topic),
		log.Int32("partition", msg.Partition),
		log.Int64("partition", msg.Offset),
		log.ByteString("key", msg.Key),
		log.ByteString("msg", msg.Value),
	)
	sess.MarkMessage(msg, "")
}

func TestKafka(t *testing.T) {
	k := New(Config{
		BootstrapServers: []string{"127.0.0.1:9092"},
	})

	k.Subscribe(ConsumerConfig{
		GroupID: "test",
		Topic:   "test",
	}, mockHandler)

	_ = k.Start(context.Background())
}
