package mid

import (
	"context"
	"github.com/ardanlabs/blockchain/foundation/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

// ctxKey represents the type of value for the context key.
// 虽然这两个 ctxKey 的基础类型都是 int，但它们是不同的类型，不可以直接相互替换。即使两个 ctxKey 类型的值看起来一样，Go 语言会认为它们是不同的，因此 ctx.Value(key) 会返回 nil，因为 key 类型不匹配。
type ctxKey int

// key is how request values are stored/retrieved.
const key = 1

func HandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建新的 Values 实例
		values := &web.Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}

		// 设置到 context.Context 中
		ctx := context.WithValue(c.Request.Context(), 1, values)
		c.Request = c.Request.WithContext(ctx)

		// 设置到 Gin 上下文中（以便后续处理中也能访问）
		c.Set("context_values", values)

		//查看是否有值
		//v, err := web.GetValues(ctx)
		//if err != nil {
		//	log.Printf("Error retrieving values from context: %v", err)
		//}
		//log.Printf("Error retrieving values from context: %v", v)
		// Continue processing the request.
		// This allows subsequent handlers to use the context values set here.
		c.Next()
	}
}
