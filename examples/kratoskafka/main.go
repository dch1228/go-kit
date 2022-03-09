package main

import (
	"context"

	"github.com/dch1228/go-kit/kafka"
	"github.com/dch1228/go-kit/log"
	"github.com/dch1228/go-kit/tracing"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Greeting struct {
	lg *log.Logger
}

func (g *Greeting) Handler(msg *kafka.Message) {
	lg := g.lg.Named("Handler").WithCtx(msg.Ctx())
	lg.Info("Greeting")
	msg.Ack()
}

func newKafkaConsumers() *kafka.Kafka {
	k := kafka.New(kafka.Config{
		BootstrapServers: []string{"127.0.0.1:9092"},
	})
	k.Use(kafka.Logger(), kafka.Recovery(), kafka.Trace(), kafka.Metric())

	greeting := &Greeting{
		lg: log.L().Named("Greeting"),
	}

	k.Subscribe(kafka.ConsumerConfig{
		GroupID: "group",
		Topic:   "topic",
	}, greeting.Handler)

	return k
}

func newMetricServer() *http.Server {
	srv := http.NewServer(
		http.Address(":8080"),
	)

	srv.Handle("/metrics", promhttp.Handler())
	srv.HandlePrefix("/debug", pprof.NewHandler())
	return srv

}

func main() {
	logger := log.New(
		log.Config{
			Level: "info",
		},
	).WithCtxFields(
		tracing.TraceID(),
		tracing.SpanID(),
	)
	log.SetLogger(logger)
	defer func() { _ = logger.Sync() }()

	if err := tracing.Init(tracing.Config{
		Name:         "kratos-kafka",
		Endpoint:     "http://127.0.0.1:14268/api/traces",
		SamplerRatio: 1,
	}); err != nil {
		panic(err)
	}
	defer func() { _ = tracing.Shutdown(context.Background()) }()

	app := kratos.New(
		kratos.Name("kratos-kafka"),
		kratos.Server(
			newKafkaConsumers(),
			newMetricServer(),
		),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
