package metrics

import (
	"nautilus/util/conf"

	"github.com/prometheus/client_golang/prometheus"
)

// prometheus
// 指标类型
//		Counter: 计数器，单调递增
//		Gauge: 计量器，值可增可减
//		Histogram: 分布图
//		Summary: 摘要
// 数据采集
//		pull:
//		push:

var (
	// RPCDurationSeconds RPC 请求耗时
	RPCDurationSeconds *prometheus.HistogramVec

	// DBDurationSeconds mysql 调用耗时
	DBDurationSeconds *prometheus.HistogramVec

	// RedisDurationSeconds redis 调用耗时
	RedisDurationSeconds *prometheus.HistogramVec

	// HTTPDurationSeconds http 调用耗时
	HTTPDurationSeconds *prometheus.HistogramVec

	// GRPCDurationSeconds grpc 调用耗时
	GRPCDurationSeconds *prometheus.HistogramVec

	// DBMaxOpenConnections 最大DB连接数
	DBMaxOpenConnections *prometheus.GaugeVec

	// DBOpenConnections 当前 DB 连接总数
	DBOpenConnections *prometheus.GaugeVec

	// DBInUseConnections 当前在用 DB 连接数
	DBInUseConnections *prometheus.GaugeVec

	// DBIdleConnections 空闲 DB 连接数
	DBIdleConnections *prometheus.GaugeVec

	// DBWaitCount 从 DB 连接池取不到连接需要等待的总数量
	DBWaitCount *prometheus.CounterVec

	// DBMaxIdleClosed 因为 SetMaxIdleConns 而被关闭的连接总数量
	DBMaxIdleClosed *prometheus.CounterVec

	// DBMaxLifetimeClosed 因为 SetMaxLifetimeClosed 而被关闭的连接总数
	DBMaxLifetimeClosed *prometheus.CounterVec
)

var buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1}

func init() {
	RPCDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "rpc_duration_seconds",
		Help:        "RPC latency distributions",
		Buckets:     buckets,
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"path", "code"})
	prometheus.MustRegister(RPCDurationSeconds)

	DBDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "db_duration_seconds",
		Help:        "MySQL latency distributions",
		Buckets:     buckets,
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name", "table", "cmd"})
	prometheus.MustRegister(DBDurationSeconds)

	RedisDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "redis_duration_seconds",
		Help:        "Redis latency distributions",
		Buckets:     buckets,
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name", "cmd"})
	prometheus.MustRegister(RedisDurationSeconds)

	HTTPDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "http_duration_seconds",
		Help:        "HTTP latency distributions",
		Buckets:     buckets,
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"url", "status"})
	prometheus.MustRegister(HTTPDurationSeconds)

	GRPCDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "grpc_duration_seconds",
		Help:        "GRPC latency distributions",
		Buckets:     buckets,
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"service", "status"})
	prometheus.MustRegister(GRPCDurationSeconds)

	DBMaxOpenConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "nautilus",
		Name:        "db_max_open_conns",
		Help:        "db max open connections",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBMaxOpenConnections)

	DBOpenConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "nautilus",
		Name:        "db_open_connections",
		Help:        "db open connections",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBOpenConnections)

	DBInUseConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "nautilus",
		Name:        "db_in_use_connections",
		Help:        "db in use connections",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBInUseConnections)

	DBIdleConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "nautilus",
		Name:        "db_idle_connections",
		Help:        "db idle connections",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBIdleConnections)

	DBWaitCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "nautilus",
		Name:        "db_wait_count",
		Help:        "db wait count",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBWaitCount)

	DBMaxIdleClosed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "nautilus",
		Name:        "db_max_idle_closed",
		Help:        "db max idle closed",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBMaxIdleClosed)

	DBMaxLifetimeClosed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "nautilus",
		Name:        "db_max_lift_time_closed",
		Help:        "db max lift time closed",
		ConstLabels: map[string]string{"app": conf.AppID},
	}, []string{"name"})
	prometheus.MustRegister(DBMaxLifetimeClosed)
}
