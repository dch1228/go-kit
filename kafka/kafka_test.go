package kafka

import (
	"context"
	"testing"
)

func mockHandler(msg *Message) {
	msg.Ack()
}

func panicHandler(msg *Message) {
	panic("panic")
	msg.Ack()
}

func TestKafka(t *testing.T) {
	k := New(Config{
		BootstrapServers: []string{"127.0.0.1:9092"},
	})

	k.Use(Logger())

	k.Subscribe(ConsumerConfig{
		GroupID: "test",
		Topic:   "test",
	}, mockHandler)

	_ = k.Start(context.Background())
}

func TestRecovery(t *testing.T) {
	k := New(Config{
		BootstrapServers: []string{"127.0.0.1:9092"},
	})

	k.Use(Recovery(), Logger())

	k.Subscribe(ConsumerConfig{
		GroupID: "test",
		Topic:   "test",
	}, panicHandler)

	_ = k.Start(context.Background())
}
