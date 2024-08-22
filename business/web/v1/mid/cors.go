package mid

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors sets the response headers needed for Cross-Origin Resource Sharing
func Cors(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the CORS headers to the response.
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If it's a preflight request, return 200 status and stop further processing.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Call the next handler.
		c.Next()
	}
}
