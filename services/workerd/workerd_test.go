package workerd

import (
	"context"
	"testing"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sourcegraph/conc"
)

func TestRunWorker(t *testing.T) {
	workerdCWD := "/home/coder/code/frp-panel/tmp/workerd"
	workerID := "test"
	workerdBinPath := "/home/coder/go/bin/workerd"

	c := context.Background()
	defaultWorker := &pb.Worker{WorkerId: &workerID}
	FillWorkerValue(defaultWorker, 1)

	if err := GenCapnpConfig(c, workerdCWD, &pb.WorkerList{Workers: []*pb.Worker{defaultWorker}}); err != nil {
		panic(err)
	}

	var wg conc.WaitGroup

	wg.Go(func() {
		time.Sleep(10 * time.Second)
	})

	if err := WriteWorkerCodeToFile(c, defaultWorker, workerdCWD); err != nil {
		panic(err)
	}

	runner := NewExecManager(workerdBinPath,
		[]string{"serve", "--watch", "--verbose"})
	runner.RunCmd(workerID, WorkerCWDPath(c, defaultWorker, workerdCWD),
		[]string{ConfigFilePath(c, defaultWorker, workerdCWD)})

	defer runner.ExitAllCmd()

	wg.Wait()
}
