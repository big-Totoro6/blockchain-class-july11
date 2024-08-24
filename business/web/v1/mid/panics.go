package mid

import (
	"fmt"
	"github.com/ardanlabs/blockchain/business/web/metrics"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
// gin自带的panic处理
//
// 功能覆盖: Gin 的 Recovery 中间件能处理 panic，记录错误，并返回 HTTP 500 响应。它的功能和你的自定义中间件类似，但省去了你编写和维护这些代码的需要。
// 扩展性: Recovery 中间件不提供直接的指标更新或其他自定义功能。如果你的应用有特定的需求，比如记录 panic 指标或特定的日志格式，你可能需要自定义中间件。
// 而我们需要自定义 因为，你可能需要自定义 panic 处理中间件的情况包括：
//
// 自定义指标: 如果你的应用需要更新特定的监控数据或指标（如 metrics.AddPanics(c)），则需要自定义中间件。
// 特殊的错误处理: 如果你需要在处理 panic 时执行特定的操作，比如自定义日志记录或发送错误报告，则可以自定义中间件。
func Panics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Defer a function to recover from a panic.
		defer func() {
			if rec := recover(); rec != nil {
				// Stack trace will be provided.
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

				// Updates the metrics stored in the context.
				metrics.AddPanics(c)

				// Log the panic or handle it according to your application's needs.
				// Optionally, you can respond with an error message to the client.
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}()

		// Proceed to the next middleware or handler.
		c.Next()
	}
}
