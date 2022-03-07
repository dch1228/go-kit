package kafka

import (
	"context"
	"testing"

	"github.com/dch1228/go-kit/tracing"
)

func mockHandler(msg *Message) {
	msg.Ack()
}

func panicHandler(_ *Message) {
	panic("panic")
}

func TestLogger(t *testing.T) {
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

	k.Use(Logger(), Recovery())

	k.Subscribe(ConsumerConfig{
		GroupID: "test",
		Topic:   "test",
	}, panicHandler)

	_ = k.Start(context.Background())
}

func TestTracing(t *testing.T) {
	if err := tracing.Init(tracing.Config{
		Name:         "TestTracing",
		Endpoint:     "http://127.0.0.1:14268/api/traces",
		SamplerRatio: 1,
	}); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tracing.Shutdown(context.Background()) }()

	k := New(Config{
		BootstrapServers: []string{"127.0.0.1:9092"},
	})

	k.Use(Logger(), Recovery(), Trace())

	k.Subscribe(ConsumerConfig{
		GroupID: "test",
		Topic:   "test",
	}, mockHandler)

	_ = k.Start(context.Background())
}
