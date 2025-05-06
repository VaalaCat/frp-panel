package workerd

import (
	"fmt"
	"strings"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/samber/lo"
)

type Opt func(*pb.Worker)

func FillWorkerValue(worker *pb.Worker, UserID uint, opt ...Opt) {

	worker.UserId = lo.ToPtr(uint32(UserID))

	if len(worker.GetName()) == 0 {
		worker.Name = lo.ToPtr(utils.NewCodeName(2))
	}

	if len(worker.GetCode()) == 0 {
		worker.Code = lo.ToPtr(string(defs.DefaultCode))
	}
	if len(worker.GetWorkerId()) == 0 {
		worker.WorkerId = lo.ToPtr(utils.GenerateUUID())
	}
	if len(worker.GetCodeEntry()) == 0 {
		worker.CodeEntry = lo.ToPtr(string(defs.DefaultEntry))
	}
	if len(worker.GetConfigTemplate()) == 0 {
		worker.ConfigTemplate = lo.ToPtr(string(defs.DefaultConfigTemplate))
	}

	worker.Socket = &pb.Socket{
		Name:    lo.ToPtr(worker.GetWorkerId()),
		Address: lo.ToPtr(fmt.Sprintf(defs.DefaultSocketTemplate, worker.GetWorkerId())),
	}

	for _, o := range opt {
		o(worker)
	}
}

func SafeWorkerID(id string) string {
	replacer := strings.NewReplacer("/", "", ".", "", "-", "")
	return replacer.Replace(id)
}
