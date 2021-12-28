package middleware

import (
	"context"
	"fmt"
	"time"

	"nautilus/pkg/ctxkit"
	"nautilus/pkg/log"
	"nautilus/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkit.StartTimeKey, time.Now())
		c.Request = c.Request.WithContext(ctx)
		defer func() {
			start := ctx.Value(ctxkit.StartTimeKey).(time.Time)
			duration := time.Since(start)
			path := c.Request.URL.Path
			status := c.Writer.Status()

			// TODO params
			log.Get(ctx).WithFields(log.Fields{
				"path":   path,
				"status": status,
				"cost":   duration.Seconds(),
			}).Info("new rpc")

			// 爬虫随便访问的url有可能导致prometheus报错
			// prometheus统计的url中不能带数字
			if status != 404 {
				metrics.RPCDurationSeconds.With(prometheus.Labels{
					"path": path,
					"code": fmt.Sprint(c.Writer.Status()),
				}).Observe(duration.Seconds())

				metrics.RPCQPSCount.With(prometheus.Labels{
					"path": path,
					"code": fmt.Sprint(c.Writer.Status()),
				}).Inc()
			}
		}()

		c.Next()
	}
}
