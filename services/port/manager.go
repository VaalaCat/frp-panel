package tunnel

import (
	"context"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

type PortManager interface {
	ClaimWorkerPort(c context.Context, workerID string) int32
	GetWorkerPort(c context.Context, workerID string) (int32, bool)
}

type portManager struct {
	portMap *utils.SyncMap[string, int32]
}

func (p *portManager) ClaimWorkerPort(c context.Context, workerID string) int32 {
	port, err := utils.GetAvailablePort(defs.DefaultHostName)
	if err != nil {
		logger.Logger(c).WithError(err).Panic("get available port failed")
	}
	p.portMap.Store(workerID, int32(port))
	return int32(port)
}

func (p *portManager) GetWorkerPort(c context.Context, workerID string) (int32, bool) {
	return p.portMap.Load(workerID)
}

var (
	mgr PortManager
)

func NewPortManager() PortManager {
	return &portManager{
		portMap: &utils.SyncMap[string, int32]{},
	}
}

func GetPortManager() PortManager {
	if mgr == nil {
		mgr = NewPortManager()
	}
	return mgr
}
