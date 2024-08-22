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
