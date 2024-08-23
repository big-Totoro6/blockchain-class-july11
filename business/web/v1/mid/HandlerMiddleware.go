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
		// 创建新的 Values 实例
		values := &web.Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}

		// 设置到 context.Context 中
		ctx := context.WithValue(c.Request.Context(), key, values)
		c.Request = c.Request.WithContext(ctx)

		// 设置到 Gin 上下文中（以便后续处理中也能访问）
		c.Set("context_values", values)

		// Continue processing the request.
		// This allows subsequent handlers to use the context values set here.
		c.Next()
	}
}
