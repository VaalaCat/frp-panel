//go:build windows

package workerd

import (
	"context"

	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

type workerExecManager struct{}

// ExitAllCmd implements app.WorkerExecManager.
func (w *workerExecManager) ExitAllCmd() {
	ctx := context.Background()
	logger.Logger(ctx).Errorf("windows has not implemented functions")
}

// ExitCmd implements app.WorkerExecManager.
func (w *workerExecManager) ExitCmd(workerId string) {
	ctx := context.Background()
	logger.Logger(ctx).Errorf("windows has not implemented functions")
}

// RunCmd implements app.WorkerExecManager.
func (w *workerExecManager) RunCmd(workerId string, cwd string, argv []string) {
	ctx := context.Background()
	logger.Logger(ctx).Errorf("windows has not implemented functions")
}

// UpdateBinaryPath implements app.WorkerExecManager.
func (w *workerExecManager) UpdateBinaryPath(path string) {
	ctx := context.Background()
	logger.Logger(ctx).Errorf("windows has not implemented functions")
}

func NewExecManager(binPath string, defaultArgs []string) app.WorkerExecManager {
	return &workerExecManager{}
}
