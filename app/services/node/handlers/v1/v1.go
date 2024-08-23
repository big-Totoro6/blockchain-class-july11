// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/ardanlabs/blockchain/app/services/node/handlers/v1/private"
	"github.com/ardanlabs/blockchain/app/services/node/handlers/v1/public"
	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
	"github.com/ardanlabs/blockchain/foundation/nameservice"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const version = "v1"

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *zap.SugaredLogger
	State *state.State
	NS    *nameservice.NameService
}

// PublicRoutes sets up the v1 routes for the given router group.
func PublicRoutes(rg *gin.RouterGroup, cfg Config) {
	pbl := public.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	}

	// Define routes and their handlers using the adapter
	rg.GET("/genesis/list", pbl.HandlerFuncAdapter(pbl.Genesis))
	rg.GET("/accounts/list", pbl.HandlerFuncAdapter(pbl.Accounts))
	rg.GET("/accounts/list/:account", pbl.HandlerFuncAdapter(pbl.Accounts))
	rg.GET("/tx/uncommitted/list", pbl.HandlerFuncAdapter(pbl.Mempool))
	rg.GET("/tx/uncommitted/list/:account", pbl.HandlerFuncAdapter(pbl.Mempool))
	rg.POST("/tx/submit", func(c *gin.Context) {
		// 直接传递 gin.Context 的 Request.Context()
		err := pbl.SubmitWalletTransaction(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})
	//rg.POST("/tx/submit", pbl.HandlerFuncAdapter(pbl.SubmitWalletTransaction))
	rg.POST("/tx/proof/:block/", pbl.HandlerFuncAdapter(pbl.SubmitWalletTransaction))
}

// PrivateRoutes binds all the version 1 private routes.
func PrivateRoutes(rg *gin.RouterGroup, cfg Config) {
	prv := private.Handlers{
		Log: cfg.Log,
	}
	rg.GET("/node/sample", prv.HandlerFuncAdapter(prv.Sample))
}
