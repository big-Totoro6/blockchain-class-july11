// Package handlers manages the different versions of the API.
package handlers

import (
	"expvar"
	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
	"github.com/ardanlabs/blockchain/foundation/nameservice"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/ardanlabs/blockchain/app/services/node/handlers/debug/checkgrp"
	v1 "github.com/ardanlabs/blockchain/app/services/node/handlers/v1"
	"github.com/ardanlabs/blockchain/business/web/v1/mid"
	"github.com/gin-contrib/cors"
	"go.uber.org/zap"
)

// MuxConfig contains all the mandatory systems required by handlers.
type MuxConfig struct {
	Shutdown chan os.Signal
	State    *state.State
	Log      *zap.SugaredLogger
	NS       *nameservice.NameService
}

// PublicMux constructs a gin.Engine with all application routes defined.
func PublicMux(cfg MuxConfig) *gin.Engine {
	// Create a new Gin engine
	r := gin.Default()

	// Apply middleware
	r.Use(mid.HandlerMiddleware())
	r.Use(mid.Logger(cfg.Log))
	r.Use(mid.Errors(cfg.Log))
	r.Use(mid.Metrics())
	//r.Use(mid.Cors("*"))
	r.Use(mid.Panics())

	//=======================跨域配置======================================================
	// 使用 cors 中间件并配置 CORS 选项
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://example.com"}, // 允许的域
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // Optional: Cache preflight responses for a day
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	//=================metrics指标访问接口=====================================================================
	// Register the /metrics route to serve expvar metrics.
	r.GET("/metrics", func(c *gin.Context) {
		expvar.Handler().ServeHTTP(c.Writer, c.Request)
	})
	//=========================================================================================

	// Handle OPTIONS requests for CORS preflight.
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204) // Respond with 'No Content' for OPTIONS preflight requests.
	})

	// Load the v1 routes
	v1Group := r.Group("/v1")
	v1.PublicRoutes(v1Group, v1.Config{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	})

	return r
}

// PrivateMux constructs a http.Handler with all application routes defined.
func PrivateMux(cfg MuxConfig) http.Handler {

	// Create a new Gin engine
	r := gin.Default()

	// Apply middleware
	r.Use(mid.HandlerMiddleware())
	r.Use(mid.Logger(cfg.Log))
	r.Use(mid.Errors(cfg.Log))
	r.Use(mid.Metrics())
	r.Use(mid.Cors("*"))
	r.Use(mid.Panics())

	// Handle OPTIONS requests for CORS preflight.
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204) // Respond with 'No Content' for OPTIONS preflight requests.
	})
	// Load the v1 routes
	v1Group := r.Group("/v1")
	v1.PrivateRoutes(v1Group, v1.Config{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	})

	return r
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}
