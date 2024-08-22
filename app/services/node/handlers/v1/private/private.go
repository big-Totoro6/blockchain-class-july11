// Package private maintains the group of handlers for node to node access.
package private

import (
	"context"
	"github.com/ardanlabs/blockchain/foundation/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// HandlerFuncAdapter 将旧的处理函数适配为 Gin 的 HandlerFunc
func (h Handlers) HandlerFuncAdapter(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := fn(c, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

// Sample just provides a starting point for the class.
func (h Handlers) Sample(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}
