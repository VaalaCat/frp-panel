package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DBManager interface {
	GetDB(dbType string) *gorm.DB
	GetDefaultDB() *gorm.DB
	SetDB(dbType string, db *gorm.DB)
	RemoveDB(dbType string)
	Init()
}

type dbManagerImpl struct {
	DBs           map[string]*gorm.DB // key: db type
	defaultDBType string
}

func (dbm *dbManagerImpl) Init() {
	for _, db := range dbm.DBs {
		if err := db.AutoMigrate(&Client{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&Client{}).TableName())
		}
		if err := db.AutoMigrate(&User{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&User{}).TableName())
		}
		if err := db.AutoMigrate(&Server{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&Server{}).TableName())
		}
		if err := db.AutoMigrate(&Cert{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&Cert{}).TableName())
		}
		if err := db.AutoMigrate(&Proxy{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&Proxy{}).TableName())
		}
		if err := db.AutoMigrate(&HistoryProxyStats{}); err != nil {
			logrus.WithError(err).Fatalf("cannot init db table [%s]", (&HistoryProxyStats{}).TableName())
		}
	}
}

var (
	dbm *dbManagerImpl
)

func MustInitDBManager(dbs map[string]*gorm.DB, defaultDBType string) {
	if dbm == nil {
		dbm = NewDBManager(dbs, defaultDBType)
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

func GetDBManager() DBManager {
	if dbm == nil {
		dbm = NewDBManager(nil, "")
	}
	return dbm
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
	return dbm.DBs[dbm.defaultDBType]
}
