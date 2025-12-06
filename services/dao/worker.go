package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
)

type WorkerQuery interface {
	GetWorkerByWorkerID(userInfo models.UserInfo, workerID string) (*models.Worker, error)
	ListWorkers(userInfo models.UserInfo, page, pageSize int) ([]*models.Worker, error)
	AdminListWorkersByClientID(clientID string) ([]*models.Worker, error)
	ListWorkersWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.Worker, error)
	CountWorkers(userInfo models.UserInfo) (int64, error)
	CountWorkersWithKeyword(userInfo models.UserInfo, keyword string) (int64, error)
}

type WorkerMutation interface {
	CreateWorker(userInfo models.UserInfo, worker *models.Worker) error
	DeleteWorker(userInfo models.UserInfo, workerID string) error
	UpdateWorker(userInfo models.UserInfo, worker *models.Worker) error
}

type workerQuery struct{ *queryImpl }
type workerMutation struct{ *mutationImpl }

func newWorkerQuery(base *queryImpl) WorkerQuery          { return &workerQuery{base} }
func newWorkerMutation(base *mutationImpl) WorkerMutation { return &workerMutation{base} }

func (m *workerMutation) CreateWorker(userInfo models.UserInfo, worker *models.Worker) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()

	worker.UserId = uint32(userInfo.GetUserID())
	worker.TenantId = uint32(userInfo.GetTenantID())

	worker.WorkerModel = nil

	return db.Create(worker).Error
}

func (m *workerMutation) DeleteWorker(userInfo models.UserInfo, workerID string) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()

	return db.Unscoped().Where(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			ID:       workerID,
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Delete(&models.Worker{}).Error
}

func (m *workerMutation) UpdateWorker(userInfo models.UserInfo, worker *models.Worker) error {
	if worker.WorkerEntity == nil {
		return fmt.Errorf("invalid worker entity")
	}
	if len(worker.WorkerEntity.ID) == 0 {
		return fmt.Errorf("invalid worker id")
	}

	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()

	if err := db.Unscoped().Model(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			ID:       worker.ID,
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Association("Clients").Unscoped().Clear(); err != nil {
		return err
	}

	return db.Where(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			ID:       worker.ID,
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Save(worker).Error
}

func (q *workerQuery) GetWorkerByWorkerID(userInfo models.UserInfo, workerID string) (*models.Worker, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	w := &models.Worker{}
	err := db.Where(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			ID:       workerID,
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Preload("Clients").First(w).Error
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (q *workerQuery) ListWorkers(userInfo models.UserInfo, page, pageSize int) ([]*models.Worker, error) {
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, fmt.Errorf("invalid page or page size")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var workers []*models.Worker
	err := db.Where(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Offset(offset).Limit(pageSize).Preload("Clients").Find(&workers).Error
	if err != nil {
		return nil, err
	}

	return workers, nil
}

func (q *workerQuery) AdminListWorkersByClientID(clientID string) ([]*models.Worker, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	client, err := newClientQuery(q.queryImpl).AdminGetClientByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = db.Model(&client).Preload("Workers").First(&client).Error
	if err != nil {
		return nil, err
	}

	return client.Workers, nil
}

func (q *workerQuery) ListWorkersWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.Worker, error) {
	if page < 1 || pageSize < 1 || len(keyword) == 0 || pageSize > 100 {
		return nil, fmt.Errorf("invalid page or page size or keyword")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var workers []*models.Worker
	err := db.Where("name like ?", "%"+keyword+"%").
		Where(&models.Worker{
			WorkerEntity: &models.WorkerEntity{
				UserId:   uint32(userInfo.GetUserID()),
				TenantId: uint32(userInfo.GetTenantID()),
			},
		}).Offset(offset).Limit(pageSize).Preload("Clients").Find(&workers).Error
	if err != nil {
		return nil, err
	}

	return workers, nil
}

func (q *workerQuery) CountWorkers(userInfo models.UserInfo) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Worker{}).Where(&models.Worker{
		WorkerEntity: &models.WorkerEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *workerQuery) CountWorkersWithKeyword(userInfo models.UserInfo, keyword string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Worker{}).Where("name like ?", "%"+keyword+"%").
		Where(&models.Worker{
			WorkerEntity: &models.WorkerEntity{
				UserId:   uint32(userInfo.GetUserID()),
				TenantId: uint32(userInfo.GetTenantID()),
			},
		}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
