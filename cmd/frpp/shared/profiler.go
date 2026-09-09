package shared

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"go.uber.org/fx"
)

func runProfiler(param struct {
	fx.In

	Lc  fx.Lifecycle
	Cfg conf.Config
	Ctx *app.Context
}) {
	if !param.Cfg.Debug.ProfilerEnabled {
		return
	}

	if !param.Cfg.IsDebug {
		logger.Logger(param.Ctx).Warn("profiler is enabled but IS_DEBUG=false, make sure you understand the risk")
	}

	addr := fmt.Sprintf(":%d", param.Cfg.Debug.ProfilerPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	param.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			logger.Logger(param.Ctx).Infof("profiler http server started: http://%s/debug/pprof/", addr)

			go func() {
				if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Logger(param.Ctx).WithError(err).Warn("profiler http server stopped unexpectedly")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
