package profile

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Enable bool
	Addr   string `default:":9090"`
}

type noop struct{}

func (*noop) Start(_ context.Context) error { return nil }
func (*noop) Stop(_ context.Context) error  { return nil }

func New(cfg Config) transport.Server {
	if !cfg.Enable {
		return new(noop)
	}
	srv := http.NewServer(
		http.Address(cfg.Addr),
	)
	srv.Handle("/metrics", promhttp.Handler())
	srv.HandlePrefix("/debug", pprof.NewHandler())
	return srv
}
