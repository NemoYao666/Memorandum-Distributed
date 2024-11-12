package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// TracingMiddleware 追踪中间件
func TracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan(c.Request.URL.Path)
		defer span.Finish()
		c.Next()
	}
}
