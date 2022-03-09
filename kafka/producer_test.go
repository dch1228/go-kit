package kafka

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel"

	"github.com/dch1228/go-kit/tracing"
)

func TestSendMessage(t *testing.T) {
	if err := tracing.Init(tracing.Config{
		Name:         "TestSendMessage",
		Endpoint:     "http://127.0.0.1:14268/api/traces",
		SamplerRatio: 1,
	}); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tracing.Shutdown(context.Background()) }()

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	producer = WrapSyncProducer(producer)

	tr := otel.Tracer("producer")
	ctx, span := tr.Start(context.Background(), "produce message")
	defer span.End()

	msg := &sarama.ProducerMessage{
		Topic: "topic",
		Key:   sarama.StringEncoder("test"),
		Value: sarama.StringEncoder("test"),
	}
	otel.GetTextMapPropagator().Inject(ctx, NewProducerMessageCarrier(msg))

	if _, _, err := producer.SendMessage(msg); err != nil {
		t.Fatal(err)
	}
}
