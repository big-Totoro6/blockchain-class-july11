package metrics

import (
	"expvar"
	"github.com/gin-gonic/gin"
	"runtime"
)

// Define a singleton instance for metrics.
var m *metrics

// metrics represents the set of metrics we gather.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// Initialize metrics.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// Set adds metrics to the Gin context.
func Set(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("metrics", m)
	c.Request = c.Request.WithContext(ctx)
}

// Update metrics functions

// AddGoroutines updates the goroutine count every 100 requests.
func AddGoroutines(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

// AddRequests increments the request count.
func AddRequests(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors increments the error count.
func AddErrors(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics increments the panic count.
func AddPanics(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		v.panics.Add(1)
	}
}
