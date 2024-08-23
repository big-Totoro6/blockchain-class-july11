package mid

import (
	"context"
	"github.com/ardanlabs/blockchain/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is how request values are stored/retrieved.
const key ctxKey = 1

func HandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中提取需要的值（如从 header 中提取 TraceID）
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String() // 如果 header 中没有 TraceID，生成一个新的
		}

		// 创建新的 Values 实例
		values := &web.Values{
			TraceID: traceID,
			Now:     time.Now().UTC(),
		}

		// 设置到 context.Context 中
		ctx := context.WithValue(c.Request.Context(), key, values)
		c.Request = c.Request.WithContext(ctx)

		// 设置到 Gin 上下文中（以便后续处理中也能访问）
		c.Set("context_values", values)

		//// 调用实际的处理函数
		//handler(c)
	}
}
