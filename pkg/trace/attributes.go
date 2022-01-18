package trace

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	// GRPCStatusCodeKey grpc status code
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
	// RPCNameKey message transmitted name
	RPCNameKey = attribute.Key("name")
	// RPCMessageTypeKey message transmitted type
	RPCMessageTypeKey = attribute.Key("message.type")
	// RPCMessageIDKey message transmitted id
	RPCMessageIDKey = attribute.Key("message.id")
	// RPCMessageCompressedSizeKey the compressed size of the message transmitted
	RPCMessageCompressedSizeKey = attribute.Key("message.compressed_size")
	// RPCMessageUncompressedSizeKey the uncompressed size of the message transmitted
	RPCMessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")

	// DBNameKey db name
	DBNameKey = semconv.DBNameKey
	// DBStatementKey sql语句
	DBStatementKey = semconv.DBStatementKey
	// DBOperationKey DML类型 select/update/insert/delete
	DBOperationKey = semconv.DBOperationKey
	// DBTableKey 表名
	DBTableKey = semconv.DBSQLTableKey
)

var (
	// RPCSystemGRPC
	RPCSystemGRPC          = semconv.RPCSystemKey.String("grpc")
	RPCNameMessage         = RPCNameKey.String("message")
	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")

	// Semantic conventions for database client calls
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/database.md#mysql

	// DBSystemValue db type
	DBSystemValue = semconv.DBSystemKey.String("mysql")
)

// StatusCodeAttr 根据给定的code返回KV
func StatusCodeAttr(code codes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(code))
}
