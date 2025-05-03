package workerd

import (
	"context"
	"errors"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func GenCapnpConfig(ctx context.Context, workerdDir string, workerList *pb.WorkerList) error {
	var hasError bool
	for _, worker := range workerList.Workers {
		fileMap := BuildCapfile([]*pb.Worker{worker})

		if fileContent, ok := fileMap[worker.GetWorkerId()]; ok {
			err := utils.WriteFile(
				ConfigFilePath(ctx, worker, workerdDir),
				fileContent)
			if err != nil {
				logrus.WithError(err).Errorf("failed to write file, worker is: %+v", worker.Name)
				hasError = true
			}
		}
	}

	logger.Logger(ctx).Infof("GenCapnpConfig has error: %v, workerList: %+v", hasError,
		lo.SliceToMap(workerList.GetWorkers(), func(w *pb.Worker) (string, bool) { return w.GetWorkerId(), true }))

	if hasError {
		return errors.New("GenCapnpConfig has error")
	}
	return nil
}
