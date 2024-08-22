// Package metrics constructs the metrics the application will track.
package metrics

import (
	"expvar"
	"github.com/gin-gonic/gin"
	"runtime"
)

// This holds the single instance of the metrics value needed for
// collecting metrics. The expvar package is already based on a singleton
// for the different metrics that are registered with the package so there
// isn't much choice here.
var m *metrics

// =============================================================================

// metrics represents the set of metrics we gather. These fields are
// safe to be accessed concurrently thanks to expvar. No extra abstraction is required.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to capture metrics.
// The metrics value is stored in a package level variable since everything
// inside of expvar is registered as a singleton. The use of once will make
// sure this initialization only happens once.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// =============================================================================

// Metrics will be supported through the context.

// ctxKeyMetric represents the type of value for the context key.
type ctxKey int

// key is how metric values are stored/retrieved.
const key ctxKey = 1

// =============================================================================

// Set sets the metrics data into the Gin context.
func Set(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("metrics", m)
	c.Request = c.Request.WithContext(ctx)
}

// Add more of these functions when a metric needs to be collected in
// different parts of the codebase. This will keep this package the
// central authority for metrics and metrics won't get lost.

// AddGoroutines refreshes the goroutine metric every 100 requests using Gin context.
func AddGoroutines(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

// AddRequests increments the request metric by 1 using Gin context.
func AddRequests(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors increments the errors metric by 1 using Gin context.
func AddErrors(c *gin.Context) {
	if v, ok := c.Value("metrics").(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics increments the panics metric by 1 using Gin context.
// Gin Context：Gin 使用 *gin.Context 来处理请求相关的上下文，而不是直接使用 context.Context。我们可以通过 c.Value("metrics_key") 来获取存储在 Gin 上下文中的 metrics 实例。
// Key管理："metrics_key" 是用于从 Gin 上下文中提取 metrics 的键。您需要确保在之前的中间件或处理器中已将 metrics 存储在 Gin 上下文中。
func AddPanics(c *gin.Context) {
	// Assuming `metrics` is stored in the Gin context with a specific key.
	if v, ok := c.Value("metrics_key").(*metrics); ok {
		v.panics.Add(1)
	}
}
