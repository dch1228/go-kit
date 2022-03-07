package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type syncProducer struct {
	tp         trace.TracerProvider
	tr         trace.Tracer
	propagator propagation.TextMapPropagator

	sarama.SyncProducer
}

func WrapSyncProducer(producer sarama.SyncProducer) sarama.SyncProducer {
	tp := otel.GetTracerProvider()
	tr := tp.Tracer("kafka.producer")

	return &syncProducer{
		tp:           tp,
		tr:           tr,
		propagator:   otel.GetTextMapPropagator(),
		SyncProducer: producer,
	}
}

func (p *syncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	span := p.startProducerSpan(msg)

	partition, offset, err = p.SyncProducer.SendMessage(msg)

	finishProducerSpan(span, partition, offset, err)

	return partition, offset, err
}

func (p *syncProducer) startProducerSpan(msg *sarama.ProducerMessage) trace.Span {
	carrier := NewProducerMessageCarrier(msg)
	ctx := p.propagator.Extract(context.Background(), carrier)

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(msg.Topic),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := p.tr.Start(ctx, "kafka.produce", opts...)
	p.propagator.Inject(ctx, carrier)

	return span
}

func finishProducerSpan(span trace.Span, partition int32, offset int64, err error) {
	span.SetAttributes(
		MessagingOffsetKey.Int64(offset),
		MessagingPartitionKey.Int64(int64(partition)),
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
