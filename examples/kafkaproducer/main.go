package main

import (
	"context"
	"flag"

	"github.com/Shopify/sarama"
	"github.com/dch1228/go-kit/conf"
	"github.com/dch1228/go-kit/kafka"
	"go.opentelemetry.io/otel"
)

var cfgPath = flag.String("c", "conf.yaml", "Specify the config file")

type Config struct {
	conf.ServerConfig
}

func main() {
	flag.Parse()
	var cfg Config

	conf.MustLoad(*cfgPath, &cfg)

	// log, trace
	cfg.MustSetup()
	defer cfg.Cleanup()

	kcfg := sarama.NewConfig()
	kcfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, kcfg)
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
