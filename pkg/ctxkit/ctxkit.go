package ctxkit

import (
	"context"
)

// key context中的key建议使用int类型
// https://github.com/golang/go/issues/17826
type key int

const (
	// TraceIDKey 请求唯一标识，类型：string
	TraceIDKey key = iota
	// StartTimeKey 请求开始时间，类型：time.Time
	StartTimeKey
	// AccessIPKey 请求的ip信息
	AccessIPKey

	// PlatformKey 平台信息，枚举 [ios, android, web, pad]
	PlatformKey
	// VersionKey 版本信息
	VersionKey
	// AccessKeyKey 登录token
	AccessKeyKey
	// AppkeyKey app key
	AppkeyKey
	// DeviceKey 浏览器型号
	DeviceKey
	// TSKey 时间戳
	TSKey
	// SignKey 签名
	SignKey
)

// GetTraceID 获取 trace id
func GetTraceID(ctx context.Context) (traceID string) {
	traceID, _ = ctx.Value(TraceIDKey).(string)
	return
}

// WithTraceID 向ctx中注入 trace id
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetAccessIP 获取请求 ip
func GetAccessIP(ctx context.Context) string {
	v, _ := ctx.Value(AccessIPKey).(string)
	return v
}

// GetPlatform 获取平台信息
func GetPlatform(ctx context.Context) string {
	v, _ := ctx.Value(PlatformKey).(string)
	return v
}

// GetAppkey 获取 app key
func GetAppkey(ctx context.Context) string {
	v, _ := ctx.Value(AppkeyKey).(string)
	return v
}

// GetDevice 获取浏览器型号
func GetDevice(ctx context.Context) string {
	v, _ := ctx.Value(DeviceKey).(string)
	return v
}
