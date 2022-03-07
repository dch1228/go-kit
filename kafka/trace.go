package kafka

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	MessagingOffsetKey    = attribute.Key("messaging.kafka.offset")
	MessagingPartitionKey = attribute.Key("messaging.kafka.partition")
)
