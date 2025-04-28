package models

import "gorm.io/gorm"

type Cert struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);uniqueIndex"`
	CertFile []byte
	KeyFile  []byte
	CaFile   []byte
}

func (c *Cert) TableName() string {
	return "certs"
}
