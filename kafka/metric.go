package kafka

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	messageTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "messages_total",
		Help: "messages_total",
	}, []string{"topic", "group_id"})

	messageDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "message_duration_millisecond",
		Help:    "message_duration_millisecond",
		Buckets: prometheus.LinearBuckets(100, 500, 5),
	}, []string{"topic", "group_id"})
)

func Metric() MiddlewareFunc {
	prometheus.MustRegister(
		messageTotal,
		messageDuration,
	)
	return func(next HandlerFunc) HandlerFunc {
		return func(msg *Message) {
			start := time.Now()

			next(msg)

			messageTotal.WithLabelValues(msg.Topic(), msg.GroupID()).Inc()
			messageDuration.WithLabelValues(msg.Topic(), msg.GroupID()).Observe(float64(time.Since(start).Milliseconds()))
		}
	}
}
