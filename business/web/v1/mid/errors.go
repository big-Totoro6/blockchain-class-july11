package mid

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/ardanlabs/blockchain/business/sys/validate"
	v1Web "github.com/ardanlabs/blockchain/business/web/v1"
	"github.com/ardanlabs/blockchain/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors and responds with a standardized format.
func Errors(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request and let other handlers execute
		c.Next()

		// Check if any errors occurred during request processing
		if len(c.Errors) > 0 {
			// Get the first error
			err := c.Errors[0].Err

			// Retrieve context values
			values, ok := c.Get("context_values")
			if !ok {
				// Handle the case where context values are missing
				log.Error("context values not found")
				// Create a shutdown error or handle it accordingly
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}

			// Type assert to *web.Values
			v, _ := values.(*web.Values)
			traceID := ""
			if v != nil {
				traceID = v.TraceID
			}

			// Log the error with trace ID
			log.Errorw("ERROR", "traceid", traceID, "ERROR", err)

			// Build out the error response
			var er v1Web.ErrorResponse
			var status int
			switch {
			case validate.IsFieldErrors(err):
				fieldErrors := validate.GetFieldErrors(err)
				er = v1Web.ErrorResponse{
					Error:  "data validation error",
					Fields: fieldErrors.Fields(),
				}
				status = http.StatusBadRequest

			case v1Web.IsRequestError(err):
				reqErr := v1Web.GetRequestError(err)
				er = v1Web.ErrorResponse{
					Error: reqErr.Error(),
				}
				status = reqErr.Status

			default:
				er = v1Web.ErrorResponse{
					Error: http.StatusText(http.StatusInternalServerError),
				}
				status = http.StatusInternalServerError
			}

			// Respond with the error to the client
			c.JSON(status, er)
		}
	}
}
