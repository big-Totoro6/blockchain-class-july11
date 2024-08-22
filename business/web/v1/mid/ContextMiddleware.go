package mid

import (
	"github.com/ardanlabs/blockchain/foundation/web"
	"github.com/gin-gonic/gin"
	"time"
)

// ContextMiddleware sets up context values for each request.
// It initializes a Values object with default or generated values,
// and stores it in the Gin context so it can be accessed later during
// the request processing lifecycle.
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new Values instance with default values.
		// Here, TraceID could be generated dynamically, e.g., using a UUID library.
		values := &web.Values{
			TraceID: "some-trace-id", // Placeholder; consider generating a unique trace ID.
			Now:     time.Now(),      // Current timestamp for the request.
		}

		// Store the Values instance in the Gin context with a specific key.
		c.Set("context_values", values)

		// Continue processing the request.
		// This allows subsequent handlers to use the context values set here.
		c.Next()
	}
}
