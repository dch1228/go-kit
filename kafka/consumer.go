package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/dch1228/go-kit/log"
)

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type HandlerFunc func(msg *Message)

type consumer struct {
	cg      sarama.ConsumerGroup
	kConfig *sarama.Config
	topics  []string
	groupID string
	handler *consumerGroupHandler
}

type consumerGroupHandler struct {
	k       *Kafka
	c       *consumer
	handler HandlerFunc
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		applyMiddleware(h.handler, h.k.middleware...)(&Message{
			msg:     msg,
			session: session,
			ctx:     context.Background(),
			groupID: h.c.groupID,
		})
	}
	return nil
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(msg *Message) {
			start := time.Now()

			next(msg)

			log.
				WithCtx(msg.ctx).
				WithOptions(log.WithCaller(false)).
				Info(
					"handle message",
					log.String("topic", msg.Topic()),
					log.Int32("partition", msg.Partition()),
					log.Int64("offset", msg.Offset()),
					log.Duration("latency", time.Since(start)),
				)
		}
	}
}

func Recovery() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(msg *Message) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					log.WithCtx(msg.ctx).Error("[PANIC RECOVER]", err)
				}
			}()
			next(msg)
		}
	}
}

func Trace() MiddlewareFunc {
	propagator := otel.GetTextMapPropagator()
	tr := otel.Tracer("kafka.consumer")
	return func(next HandlerFunc) HandlerFunc {
		return func(msg *Message) {
			parentSpanCtx := propagator.Extract(msg.ctx, msg)

			attrs := []attribute.KeyValue{
				semconv.MessagingSystemKey.String("kafka"),
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingDestinationKey.String(msg.Topic()),
				semconv.MessagingOperationReceive,
				MessagingPartitionKey.Int64(int64(msg.Partition())),
				MessagingOffsetKey.Int64(msg.Offset()),
			}
			opts := []trace.SpanStartOption{
				trace.WithAttributes(attrs...),
				trace.WithSpanKind(trace.SpanKindConsumer),
			}

			newCtx, span := tr.Start(parentSpanCtx, "kafka.consume", opts...)
			propagator.Inject(newCtx, msg)
			defer span.End()

			next(msg)
		}
	}
}
