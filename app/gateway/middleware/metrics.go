package middleware

import (
	"github.com/gin-gonic/gin"
	"micro-todoList-k8s/app/gateway/metrics"
	"micro-todoList-k8s/config"
	"strconv"
	"time"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		metrics.GateWayRequestCounter.WithLabelValues(
			config.C.Server.GateWayServiceName,
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		metrics.GateWayRequestDuration.WithLabelValues(
			config.C.Server.GateWayServiceName,
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
