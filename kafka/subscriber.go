package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"

	"github.com/dch1228/go-kit/log"
)

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type HandlerFunc func(msg *Message)

type subscriber struct {
	cg      sarama.ConsumerGroup
	kConfig *sarama.Config
	topics  []string
	groupID string
	handler *consumerGroupHandler
}

type consumerGroupHandler struct {
	k       *Kafka
	handler HandlerFunc
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		applyMiddleware(h.handler, h.k.middleware...)(&Message{
			msg:     msg,
			session: session,
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

			log.WithCtx(msg.ctx).Info("handle message",
				log.String("topic", msg.Topic()),
				log.Int32("partition", msg.Partition()),
				log.Int64("offset", msg.Offset()),
				log.ByteString("key", msg.Key()),
				log.ByteString("msg", msg.Value()),
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
	return func(next HandlerFunc) HandlerFunc {
		return func(msg *Message) {
			// todo
		}
	}
}
