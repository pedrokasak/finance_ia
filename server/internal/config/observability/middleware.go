package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware records request metrics for every route
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		RequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		RequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}
