package workerd

import (
	"context"
	"path/filepath"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
)

func WriteWorkerCodeToFile(ctx context.Context, worker *pb.Worker, workerdCWD string) error {
	return utils.WriteFile(
		CodeFilePath(ctx, worker, workerdCWD),
		string(worker.GetCode()))
}

func CodeFilePath(ctx context.Context, worker *pb.Worker, workerdCWD string) string {
	return filepath.Join(
		WorkerCWDPath(ctx, worker, workerdCWD),
		defs.WorkerCodePath,
		worker.GetCodeEntry())
}

func WorkerCodeRootPath(ctx context.Context, worker *pb.Worker, workerdCWD string) string {
	return filepath.Join(
		WorkerCWDPath(ctx, worker, workerdCWD),
		defs.WorkerCodePath)
}

func WorkerCWDPath(ctx context.Context, worker *pb.Worker, workerdCWD string) string {
	return filepath.Join(
		workerdCWD,
		defs.WorkerInfoPath,
		worker.GetWorkerId(),
	)
}

func ConfigFilePath(ctx context.Context, worker *pb.Worker, workerdCWD string) string {
	return filepath.Join(
		WorkerCWDPath(ctx, worker, workerdCWD),
		defs.CapFileName,
	)
}
