package workerd

import (
	"context"
	"os"
	"strings"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

var _ app.WorkerController = (*workerdController)(nil)

type workerdController struct {
	worker     *pb.Worker
	workerdCwd string
	status     *defs.WorkerStatus
}

func NewWorkerdController(worker *pb.Worker, workerdCwd string) *workerdController {
	return &workerdController{
		worker:     worker,
		workerdCwd: workerdCwd,
	}
}

func (w *workerdController) RunWorker(c *app.Context) {
	if err := w.Init(c); err != nil {
		logger.Logger(c).WithError(err).Errorf("init worker failed, workerId: [%s]", w.worker.GetWorkerId())
		return
	}

	execMgr := c.GetApp().GetWorkerExecManager()
	execMgr.RunCmd(
		w.worker.GetWorkerId(), WorkerCWDPath(c, w.worker, w.workerdCwd),
		[]string{ConfigFilePath(c, w.worker, w.workerdCwd)},
	)
}

func (w *workerdController) StopWorker(c *app.Context) {
	execMgr := c.GetApp().GetWorkerExecManager()
	execMgr.ExitCmd(w.worker.GetWorkerId())
	w.GarbageCollect()
}

func (w *workerdController) GetWorkerStatus(c *app.Context) defs.WorkerStatus {
	if w.status == nil {
		return defs.WorkerStatus_Unknown
	}
	return *w.status
}

func (w *workerdController) Init(c *app.Context) error {
	workerCodePath := WorkerCodeRootPath(c, w.worker, w.workerdCwd)

	// 1. 创建工作目录
	if err := os.MkdirAll(workerCodePath, os.ModePerm); err != nil {
		logger.Logger(c).WithError(err).Errorf("create work dir failed, path: [%s]", workerCodePath)
		return err
	}

	// 2. 写入配置文件和代码文件
	if err := WriteWorkerCodeToFile(c, w.worker, w.workerdCwd); err != nil {
		logger.Logger(c).WithError(err).Errorf("write worker code failed, workerId: [%s]", w.worker.GetWorkerId())
		return err
	}

	if err := GenCapnpConfig(c, w.workerdCwd, &pb.WorkerList{Workers: []*pb.Worker{w.worker}}); err != nil {
		logger.Logger(c).WithError(err).Errorf("gen worker capnp config failed, workerId: [%s]", w.worker.GetWorkerId())
		return err
	}

	logger.Logger(c).Infof("init worker success, workerId: [%s], code path: [%s]", w.worker.GetWorkerId(), workerCodePath)

	return nil
}

func (w *workerdController) GarbageCollect() {
	ctx := context.Background()

	pathToRemove := WorkerCWDPath(ctx, w.worker, w.workerdCwd)

	if !strings.HasPrefix(pathToRemove, "/tmp") {
		logger.Logger(ctx).Errorf("path not start with /tmp, do not remove path: [%s]", pathToRemove)
		return
	}

	if err := os.RemoveAll(pathToRemove); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("remove path failed, path: [%s]", pathToRemove)
	}
}
