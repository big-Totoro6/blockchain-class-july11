package mid

import (
	"github.com/ardanlabs/blockchain/foundation/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Values from the context.
		values, ok := c.Get("context_values")
		if !ok {
			// Handle the missing context values case (log error or return response).
			log.Error("context values not found")
			c.JSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Type assert to *web.Values.
		v, _ := values.(*web.Values)

		// Log request start.
		start := time.Now()
		log.Infow("request started", "traceid", v.TraceID, "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.ClientIP())

		// Process the request.
		c.Next()

		// Log request completion.
		log.Infow("request completed", "traceid", v.TraceID, "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.ClientIP(), "statuscode", c.Writer.Status(), "latency", time.Since(start))
	}
}
