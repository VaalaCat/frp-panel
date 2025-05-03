package workerd

import (
	"fmt"
	"runtime"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

type workersManager struct {
	workers *utils.SyncMap[string, app.WorkerController]
}

func NewWorkersManager() *workersManager {
	return &workersManager{
		workers: &utils.SyncMap[string, app.WorkerController]{},
	}
}

func (m *workersManager) GetWorker(ctx *app.Context, id string) (app.WorkerController, bool) {
	return m.workers.Load(id)
}

func (m *workersManager) RunWorker(ctx *app.Context, id string, worker app.WorkerController) error {
	if !ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		logger.Logger(ctx).Errorf("function features are not enabled")
		return fmt.Errorf("function features are not enabled")
	}

	worker.RunWorker(ctx)

	m.workers.Store(id, worker)
	return nil
}

func (m *workersManager) StopWorker(ctx *app.Context, id string) error {
	worker, ok := m.workers.Load(id)
	if !ok {
		return fmt.Errorf("cannot find worker, id: %s", id)
	}
	worker.StopWorker(ctx)
	m.workers.Delete(id)
	return nil
}

func (m *workersManager) StopAll() {
	m.workers.Range(func(k string, v app.WorkerController) bool {
		v.StopWorker(nil)
		return true
	})

	tmpM := m.workers.ToMap()
	for k := range tmpM {
		m.workers.Delete(k)
	}
}

func (m *workersManager) GetWorkerStatus(ctx *app.Context, id string) (defs.WorkerStatus, error) {
	ok, err := utils.ProcessExistsBySelf(id)
	if err != nil {
		return defs.WorkerStatus_Unknown, err
	}
	if ok {
		return defs.WorkerStatus_Running, nil
	}
	return defs.WorkerStatus_Inactive, nil
}

func (m *workersManager) InstallWorkerd(ctx *app.Context, url string, installDir string) (string, error) {
	arch := runtime.GOARCH
	os := runtime.GOOS

	workerDownloadCfg := ctx.GetApp().GetConfig().Client.Worker.WorkerdDownloadURL

	if os != "linux" {
		return "", fmt.Errorf("unsupported os: %s", os)
	}
	if arch != "amd64" && arch != "arm64" {
		return "", fmt.Errorf("unsupported arch: %s", arch)
	}

	downloadUrl := ""
	if len(url) > 0 {
		downloadUrl = url
	} else {
		switch arch {
		case "amd64":
			downloadUrl = workerDownloadCfg.LinuxX8664
		case "arm64":
			downloadUrl = workerDownloadCfg.LinuxArm64
		default:
			return "", fmt.Errorf("unsupported arch: %s", arch)
		}
	}

	if workerDownloadCfg.UseProxy {
		if len(ctx.GetApp().GetConfig().App.GithubProxyUrl) > 0 {
			downloadUrl = fmt.Sprintf("%s/%s", ctx.GetApp().GetConfig().App.GithubProxyUrl, downloadUrl)
		}
	}

	proxyUrl := ctx.GetApp().GetConfig().HTTP_PROXY

	path, err := utils.DownloadFile(ctx, downloadUrl, proxyUrl)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("failed to download workerd, url: %s", downloadUrl)
		return "", err
	}

	if len(installDir) == 0 {
		installDir = "/usr/local/bin"
	}

	finalPath, err := utils.ExtractGZTo(path, "workerd", installDir)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("failed to extract workerd, path: %s", path)
		return "", err
	}

	logger.Logger(ctx).Infof("workerd installed successfully, path: %s", finalPath)

	return finalPath, nil
}
