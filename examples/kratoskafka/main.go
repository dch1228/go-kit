package main

import (
	"flag"

	"github.com/dch1228/go-kit/conf"
	"github.com/dch1228/go-kit/kafka"
	"github.com/dch1228/go-kit/log"
	"github.com/dch1228/go-kit/profile"
	"github.com/go-kratos/kratos/v2"
)

var cfgPath = flag.String("c", "conf.yaml", "Specify the config file")

type Config struct {
	conf.ServerConfig
}

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

func main() {
	flag.Parse()
	var cfg Config

	conf.MustLoad(*cfgPath, &cfg)

	// log, trace
	cfg.MustSetup()
	defer cfg.Cleanup()

	log.Info("config", log.Any("config", cfg))

	app := kratos.New(
		kratos.Name("kratos-kafka"),
		kratos.Server(
			newKafkaConsumers(),
			profile.New(cfg.Profile),
		),
	)

	if err := app.Run(); err != nil {
		log.Error("app shutdown", err)
	} else {
		log.Info("app shutdown")
	}
}
