//go:build !windows
// +build !windows

package wg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/services/wg/multibind"
	"github.com/VaalaCat/frp-panel/services/wg/transport/ws"

	"golang.zx2c4.com/wireguard/conn"
)

func (w *wireGuard) initTransports() error {
	log := w.svcLogger.WithField("op", "initTransports")

	wsTrans := ws.NewWSBind(w.ctx)
	w.multiBind = multibind.NewMultiBind(
		w.svcLogger,
		multibind.NewTransport(conn.NewDefaultBind(), "udp"),
		multibind.NewTransport(wsTrans, "ws"),
	)

	engine := gin.New()
	engine.Any(defs.DefaultWSHandlerPath, func(c *gin.Context) {
		err := wsTrans.HandleHTTP(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	// if ws listen port not set, use wg listen port, share tcp and udp port
	listenPort := w.ifce.GetWsListenPort()
	if listenPort == 0 {
		listenPort = w.ifce.GetListenPort()
	}
	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", listenPort)); err != nil {
			w.svcLogger.WithError(err).Errorf("failed to run gin engine for ws transport on port %d", listenPort)
		}
	}()

	log.Infof("WS transport engine running on port %d", listenPort)

	return nil
}
