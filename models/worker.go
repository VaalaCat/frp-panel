package models

import (
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type WorkerModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Worker struct {
	*WorkerModel
	*WorkerEntity
	Clients []Client `gorm:"many2many:worker_clients;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WorkerEntity struct {
	ID             string `gorm:"type:varchar(255);uniqueIndex;not null;primaryKey"`
	Name           string `gorm:"type:varchar(255);index"`
	UserId         uint32 `gorm:"index"`
	TenantId       uint32 `gorm:"index"`
	Socket         JSON[*pb.Socket]
	CodeEntry      string
	Code           string
	ConfigTemplate string
}

func (w *Worker) TableName() string {
	return "workers"
}

func (w *WorkerEntity) FromPB(worker *pb.Worker) *WorkerEntity {
	w.ID = worker.GetWorkerId()
	w.Name = worker.GetName()
	w.UserId = uint32(worker.GetUserId())
	w.TenantId = uint32(worker.GetTenantId())
	w.Socket = JSON[*pb.Socket]{Data: worker.GetSocket()}
	w.CodeEntry = worker.GetCodeEntry()
	w.Code = worker.GetCode()
	w.ConfigTemplate = worker.GetConfigTemplate()

	return w
}

func (w *WorkerEntity) ToPB() *pb.Worker {
	return &pb.Worker{
		WorkerId:       lo.ToPtr(w.ID),
		Name:           lo.ToPtr(w.Name),
		UserId:         lo.ToPtr(uint32(w.UserId)),
		TenantId:       lo.ToPtr(uint32(w.TenantId)),
		Socket:         w.Socket.Data,
		CodeEntry:      lo.ToPtr(w.CodeEntry),
		Code:           lo.ToPtr(w.Code),
		ConfigTemplate: lo.ToPtr(w.ConfigTemplate),
	}
}

func (w *Worker) FromPB(worker *pb.Worker) *Worker {
	if w.WorkerEntity == nil {
		w.WorkerEntity = &WorkerEntity{}
	}
	if w.WorkerModel == nil {
		w.WorkerModel = &WorkerModel{}
	}

	w.WorkerEntity = w.WorkerEntity.FromPB(worker)
	return w
}

func (w *Worker) ToPB() *pb.Worker {
	if w.WorkerEntity == nil {
		w.WorkerEntity = &WorkerEntity{}
	}
	if w.WorkerModel == nil {
		w.WorkerModel = &WorkerModel{}
	}

	ret := w.WorkerEntity.ToPB()
	return ret
}
