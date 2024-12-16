package models

import (
	"time"

	"gorm.io/gorm"
)

type ProxyStats struct {
	*ProxyStatsEntity
}

type ProxyStatsEntity struct {
	ProxyID           int    `json:"proxy_id" gorm:"primary_key;auto_increment"`
	ServerID          string `json:"server_id" gorm:"index"`
	ClientID          string `json:"client_id" gorm:"index"`
	OriginClientID    string `json:"origin_client_id" gorm:"index"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	UserID            int    `json:"user_id" gorm:"index"`
	TenantID          int    `json:"tenant_id" gorm:"index"`
	TodayTrafficIn    int64  `json:"today_traffic_in"`
	TodayTrafficOut   int64  `json:"today_traffic_out"`
	HistoryTrafficIn  int64  `json:"history_traffic_in"`
	HistoryTrafficOut int64  `json:"history_traffic_out"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (*ProxyStats) TableName() string {
	return "proxy_stats"
}

// HistoryProxyStats 历史流量统计，不保证精准，只是为了展示。精准请使用 proxies 表中的 history_traffic_in/history_traffic_out
// 后续看看是否要改成时序类数据 https://github.com/nakabonne/tstorage
type HistoryProxyStats struct {
	gorm.Model
	ProxyID        int    `json:"proxy_id" gorm:"index"`
	ServerID       string `json:"server_id" gorm:"index"`
	ClientID       string `json:"client_id" gorm:"index"`
	OriginClientID string `json:"origin_client_id" gorm:"index"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	UserID         int    `json:"user_id" gorm:"index"`
	TenantID       int    `json:"tenant_id" gorm:"index"`
	TrafficIn      int64  `json:"traffic_in"`
	TrafficOut     int64  `json:"traffic_out"`
}

func (*HistoryProxyStats) TableName() string {
	return "history_proxy_stats"
}
