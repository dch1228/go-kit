package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type Message struct {
	ctx     context.Context
	msg     *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
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
