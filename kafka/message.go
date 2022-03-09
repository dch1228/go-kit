package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type Message struct {
	ctx     context.Context
	msg     *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
	groupID string
}

type ProducerMessageCarrier struct {
	msg *sarama.ProducerMessage
}

func (msg *Message) Ack() {
	msg.session.MarkMessage(msg.msg, "")
}

func (msg *Message) Topic() string {
	return msg.msg.Topic
}

func (msg *Message) Partition() int32 {
	return msg.msg.Partition
}

func (msg *Message) Offset() int64 {
	return msg.msg.Offset
}

func (msg *Message) Key() []byte {
	return msg.msg.Key
}

func (msg *Message) Value() []byte {
	return msg.msg.Value
}

func (msg *Message) GroupID() string {
	return msg.groupID
}

func (msg *Message) Ctx() context.Context {
	return msg.ctx
}

func (msg *Message) Get(key string) string {
	for _, h := range msg.msg.Headers {
		if h != nil && string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (msg *Message) Set(key, val string) {
	// 删掉重复的 key
	for i := 0; i < len(msg.msg.Headers); i++ {
		if msg.msg.Headers[i] != nil && string(msg.msg.Headers[i].Key) == key {
			msg.msg.Headers = append(msg.msg.Headers[:i], msg.msg.Headers[i+1:]...)
			i--
		}
	}
	msg.msg.Headers = append(msg.msg.Headers, &sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(val),
	})
}

func (msg *Message) Keys() []string {
	out := make([]string, len(msg.msg.Headers))
	for i, h := range msg.msg.Headers {
		out[i] = string(h.Key)
	}
	return out
}

func NewProducerMessageCarrier(msg *sarama.ProducerMessage) ProducerMessageCarrier {
	return ProducerMessageCarrier{msg: msg}
}

func (c ProducerMessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c ProducerMessageCarrier) Set(key, val string) {
	for i := 0; i < len(c.msg.Headers); i++ {
		if string(c.msg.Headers[i].Key) == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
			i--
		}
	}
	c.msg.Headers = append(c.msg.Headers, sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(val),
	})
}

func (c ProducerMessageCarrier) Keys() []string {
	out := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		out[i] = string(h.Key)
	}
	return out
}
