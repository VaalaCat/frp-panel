package models

import (
	"time"

	"gorm.io/gorm"
)

type Proxy struct {
	*ProxyEntity
}

type ProxyEntity struct {
	ProxyID           int    `json:"proxy_id" gorm:"primary_key;auto_increment"`
	ServerID          string `json:"server_id" gorm:"index"`
	ClientID          string `json:"client_id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	UserID            int    `json:"user_id"`
	TenantID          int    `json:"tenant_id"`
	TodayTrafficIn    int64  `json:"today_traffic_in"`
	TodayTrafficOut   int64  `json:"today_traffic_out"`
	HistoryTrafficIn  int64  `json:"history_traffic_in"`
	HistoryTrafficOut int64  `json:"history_traffic_out"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (*Proxy) TableName() string {
	return "proxies"
}
