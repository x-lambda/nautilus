package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

const messageEvent = "message"

var (
	// MessageSent type of sent message
	MessageSent = messageType(RPCMessageTypeSent)
	// MessageReceived type of received message
	MessageReceived = messageType(RPCMessageTypeReceived)
)

type messageType attribute.KeyValue

func (m messageType) Event(ctx context.Context, id int, message interface{}) {
	span := trace.SpanFromContext(ctx)

	if p, ok := message.(proto.Message); ok {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
			RPCMessageUncompressedSizeKey.Int(proto.Size(p)),
		))
	} else {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
		))
	}
}
