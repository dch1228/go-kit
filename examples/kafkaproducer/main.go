package main

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/dch1228/go-kit/kafka"
	"github.com/dch1228/go-kit/tracing"
	"go.opentelemetry.io/otel"
)

func main() {
	if err := tracing.Init(tracing.Config{
		Name:         "kafka-producer",
		Endpoint:     "http://127.0.0.1:14268/api/traces",
		SamplerRatio: 1,
	}); err != nil {
		panic(err)
	}
	defer func() { _ = tracing.Shutdown(context.Background()) }()

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, cfg)
	if err != nil {
		panic(err)
	}
	producer = kafka.WrapSyncProducer(producer)

	tr := otel.Tracer("producer")
	ctx, span := tr.Start(context.Background(), "produce message")
	defer span.End()

	msg := &sarama.ProducerMessage{
		Topic: "topic",
		Key:   sarama.StringEncoder("test"),
		Value: sarama.StringEncoder("test"),
	}
	otel.GetTextMapPropagator().Inject(ctx, kafka.NewProducerMessageCarrier(msg))

	if _, _, err := producer.SendMessage(msg); err != nil {
		panic(err)
	}
}
