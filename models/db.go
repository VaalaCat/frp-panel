package models

import (
	"context"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"gorm.io/gorm"
)

type dbManagerImpl struct {
	DBs           map[string]map[string]*gorm.DB // map[db type]map[db role]*gorm.DB
	defaultDBType string
	debug         bool
}

func (dbm *dbManagerImpl) Init() {
	for _, dbGroup := range dbm.DBs {
		for _, db := range dbGroup {
			if err := db.AutoMigrate(&Client{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&Client{}).TableName())
			}
			if err := db.AutoMigrate(&User{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&User{}).TableName())
			}
			if err := db.AutoMigrate(&Server{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&Server{}).TableName())
			}
			if err := db.AutoMigrate(&Cert{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&Cert{}).TableName())
			}
			if err := db.AutoMigrate(&ProxyStats{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&ProxyStats{}).TableName())
			}
			if err := db.AutoMigrate(&HistoryProxyStats{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&HistoryProxyStats{}).TableName())
			}
			if err := db.AutoMigrate(&Worker{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&Worker{}).TableName())
			}
			if err := db.AutoMigrate(&ProxyConfig{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&ProxyConfig{}).TableName())
			}
			if err := db.AutoMigrate(&UserGroup{}); err != nil {
				logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&UserGroup{}).TableName())
			}
		}
	}
}

func NewDBManager(defaultDBType string) *dbManagerImpl {
	dbs := map[string]map[string]*gorm.DB{}
	return &dbManagerImpl{
		DBs:           dbs,
		defaultDBType: defaultDBType,
	}
}

func (dbm *dbManagerImpl) GetDB(dbType string, dbRole string) *gorm.DB {
	return dbm.DBs[dbType][dbRole]
}

func (dbm *dbManagerImpl) SetDB(dbType string, dbRole string, db *gorm.DB) {
	if dbm.DBs[dbType] == nil {
		dbm.DBs[dbType] = map[string]*gorm.DB{}
	}
	dbm.DBs[dbType][dbRole] = db
}

func (dbm *dbManagerImpl) RemoveDB(dbType string, dbRole string) {
	if dbm.DBs[dbType] == nil {
		return
	}
	delete(dbm.DBs[dbType], dbRole)
}

func (dbm *dbManagerImpl) GetDefaultDB() *gorm.DB {
	dbGroup := dbm.DBs[dbm.defaultDBType]
	if dbm.debug {
		return dbGroup[defs.DBRoleDefault].Debug()
	}
	return dbGroup[defs.DBRoleDefault]
}

func (dbm *dbManagerImpl) SetDebug(debug bool) {
	dbm.debug = debug
}
