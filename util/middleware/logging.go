package middleware

import (
	"context"
	"time"

	"nautilus/util/ctxkit"
	"nautilus/util/log"

	"github.com/gin-gonic/gin"
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

			// TODO params
			log.Get(ctx).WithFields(log.Fields{
				"path":   path,
				"status": c.Writer.Status(),
				"cost":   duration.Seconds(),
			}).Info("new rpc")
		}()

		c.Next()
	}
}
