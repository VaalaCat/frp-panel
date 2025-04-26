package models

import (
	"context"

	"github.com/VaalaCat/frp-panel/logger"
	"gorm.io/gorm"
)

type DBManager interface {
	GetDB(dbType string) *gorm.DB
	GetDefaultDB() *gorm.DB
	SetDB(dbType string, db *gorm.DB)
	RemoveDB(dbType string)
	SetDebug(bool)
	Init()
}

type dbManagerImpl struct {
	DBs           map[string]*gorm.DB // key: db type
	defaultDBType string
	debug         bool
}

func (dbm *dbManagerImpl) Init() {
	for _, db := range dbm.DBs {
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
		if err := db.AutoMigrate(&ProxyConfig{}); err != nil {
			logger.Logger(context.Background()).WithError(err).Fatalf("cannot init db table [%s]", (&ProxyConfig{}).TableName())
		}
	}
}

func NewDBManager(dbs map[string]*gorm.DB, defaultDBType string) *dbManagerImpl {
	if dbs == nil {
		dbs = make(map[string]*gorm.DB)
	}
	return &dbManagerImpl{
		DBs:           dbs,
		defaultDBType: defaultDBType,
	}
}

func (dbm *dbManagerImpl) GetDB(dbType string) *gorm.DB {
	return dbm.DBs[dbType]
}

func (dbm *dbManagerImpl) SetDB(dbType string, db *gorm.DB) {
	dbm.DBs[dbType] = db
}

func (dbm *dbManagerImpl) RemoveDB(dbType string) {
	delete(dbm.DBs, dbType)
}

func (dbm *dbManagerImpl) GetDefaultDB() *gorm.DB {
	db := dbm.DBs[dbm.defaultDBType]
	if dbm.debug {
		return db.Debug()
	}
	return db
}

func (dbm *dbManagerImpl) SetDebug(debug bool) {
	dbm.debug = debug
}
