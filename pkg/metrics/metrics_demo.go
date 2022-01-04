package metrics

import (
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
)

/**
 * 1. Metrics
 * 2. Custom Collectors and constant Metrics
 * 3. Advanced Uses of the Registry
 * 4. HTTP Exposition
 * 5. Pushing to the PushGateway
 * 6. Graphite Bridge
 * 7. Other Means of Exposition
 *
 * prometheus client采集指标时，通过使用atomic来避免锁竞争
 *	atomic是CPU级别指令控制，目标是避免多个线程同时操作一个数值，针对的是具体的数值级别的
 *	mutex是通过一系列状态值来避免竞争，目标是保证一段代码区间不重入，侧重点是一段代码区间
 *
 * prometheus自带监控指标
 * https://docs.rancher.cn/docs/octopus/monitoring/_index/#prometheus-%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8C%87%E6%A0%87%E5%AF%B9%E7%85%A7%E8%A1%A8
 */

var (
	// CPUTemp: Gauge 仪表盘
	// 代表一种样本数据可以任意变化的指标，即可增可减
	// 程序中可以使用的地方：DB连接池的总数 当前活跃的DB连接数
	CPUTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU",
	})

	// HDFailures CounterVec 计数器
	// 代表一种样本数据单调递增的指标，只增不减，可以通过内置的函数展示事件产生的速率的变化
	// 程序中可以使用的地方 HTTP QPS
	HDFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors.",
	}, []string{"device"})

	// RPCQPSCountDemo HTTP QPS 统计
	// 需要通过内置函数来统计QPS
	// sum(rate(nautilus_rpc_qps_count [1m])) by (path)
	RPCQPSCountDemo = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "nautilus",
		Name:        "rpc_qps_count",
		Help:        "RPC QPS Count",
		ConstLabels: map[string]string{"app_id": "nautilus.server", "env": "prod"},
	}, []string{"path", "code"})

	// RPCDurationSecondsDemo 统计http接口耗时
	// Histogram 直方图
	// 某些量化指标的平均值，比如小于200ms的接口占比，200ms-500ms的接口占比，大于500ms的接口占比
	// 程序中可以用来统计http接口/redis/db的耗时
	RPCDurationSecondsDemo = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "nautilus",
		Name:        "rpc_duration_seconds",
		Help:        "RPC latency distributions",
		Buckets:     []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		ConstLabels: map[string]string{"app_id": "nautilus.server", "env": "prod"},
	}, []string{"path", "code"})
)

// _demo 只是个demo
func _demo() {
	prometheus.MustRegister(CPUTemp)
	prometheus.MustRegister(HDFailures)
	prometheus.MustRegister(RPCQPSCountDemo)
	prometheus.MustRegister(RPCDurationSecondsDemo)

	// nautilus_rpc_qps_count{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order"} 1
	// nautilus_rpc_qps_count{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile"} 1
	RPCQPSCount.With(prometheus.Labels{"path": "/api/user.v1.info/profile", "code": "200"}).Inc()
	RPCQPSCount.With(prometheus.Labels{"path": "/api/user.v1.info/order", "code": "200"}).Inc()

	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.005"} 7
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.01"} 14
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.025"} 25
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.05"} 46
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.1"} 105
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.25"} 264
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="0.5"} 525
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="1"} 1000
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile",le="+Inf"} 1000
	//nautilus_rpc_duration_seconds_sum{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile"} 480.7329999999993
	//nautilus_rpc_duration_seconds_count{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/profile"} 1000
	for i := 0; i < 1000; i++ {
		// 1000ms之内
		x := (float64(rand.Int31() % 1000)) / 1000
		RPCDurationSecondsDemo.With(prometheus.Labels{"path": "/api/user.v1.info/profile", "code": "200"}).Observe(x)
	}

	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.005"} 17   落在0.005秒之内的请求有17个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.01"} 26    落在0.01 秒之内的请求有26个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.025"} 48   落在0.025秒之内的请求有48个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.05"} 84    落在0.05 秒之内的请求有84个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.1"} 171    落在0.1  秒之内的请求有171个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.25"} 492   落在0.25 秒之内的请求有492个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="0.5"} 1000   落在0.5  秒之内的请求有1000个
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="1"} 1000
	//nautilus_rpc_duration_seconds_bucket{app_id="nautilus.server",code="200",env="prod",path="/api/user.v1.info/order",le="+Inf"} 1000
	for i := 0; i < 1000; i++ {
		// 500ms之内
		x := (float64(rand.Int31() % 500)) / 1000
		RPCDurationSecondsDemo.With(prometheus.Labels{"path": "/api/user.v1.info/order", "code": "200"}).Observe(x)
	}
}
