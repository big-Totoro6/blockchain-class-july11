package mid

import (
	"github.com/ardanlabs/blockchain/business/web/metrics"
	"github.com/gin-gonic/gin"
)

// Metrics 适配步骤
// 使用 Gin 上下文：用 *gin.Context 替代 context.Context。
// 设置和更新指标：将指标相关的操作适配到 Gin 上下文中，并在请求处理过程中更新这些指标。
// 记录指标数据：在 Gin 上下文中设置和更新指标，确保所有中间件和处理程序都可以访问和更新这些指标。
// Metrics updates program counters.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set metrics into the Gin context.
		metrics.Set(c)

		// Process the request.
		c.Next()

		// After the request has been processed, update the metrics.
		metrics.AddRequests(c)
		metrics.AddGoroutines(c)

		// Check if there were any errors and update the metrics accordingly.
		if len(c.Errors) > 0 {
			metrics.AddErrors(c)
		}
	}
}
