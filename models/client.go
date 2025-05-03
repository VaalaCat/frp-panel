package models

import (
	"encoding/json"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Client struct {
	*ClientEntity
	Workers []*Worker `json:"workers,omitempty" gorm:"many2many:worker_clients;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ClientEntity struct {
	ClientID       string `json:"client_id" gorm:"uniqueIndex;not null;primaryKey"`
	ServerID       string `json:"server_id"`
	TenantID       int    `json:"tenant_id" gorm:"not null"`
	UserID         int    `json:"user_id" gorm:"not null"`
	ConfigContent  []byte `json:"config_content"`
	ConnectSecret  string `json:"connect_secret" gorm:"not null"`
	Stopped        bool   `json:"stopped"`
	Comment        string `json:"comment"`
	IsShadow       bool   `json:"is_shadow" gorm:"index"`
	OriginClientID string `json:"origin_client_id" gorm:"index"`
	FrpsUrl        string `json:"frps_url" gorm:"index"`
	Ephemeral      bool   `json:"ephemeral" gorm:"index"`

	LastSeenAt *time.Time `json:"last_seen_at" gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (*Client) TableName() string {
	return "clients"
}

func (c *ClientEntity) SetConfigContent(cfg v1.ClientConfig) error {
	newCfg := struct {
		v1.ClientCommonConfig
		Proxies  []v1.ProxyConfigurer   `json:"proxies,omitempty"`
		Visitors []v1.VisitorBaseConfig `json:"visitors,omitempty"`
	}{
		ClientCommonConfig: cfg.ClientCommonConfig,
		Proxies: lo.Map(cfg.Proxies, func(item v1.TypedProxyConfig, _ int) v1.ProxyConfigurer {
			return item.ProxyConfigurer
		}),
		Visitors: lo.Map(cfg.Visitors, func(item v1.TypedVisitorConfig, _ int) v1.VisitorBaseConfig {
			return *item.GetBaseConfig()
		}),
	}
	raw, err := json.Marshal(newCfg)
	if err != nil {
		return err
	}
	c.ConfigContent = raw
	return nil
}

func (c *ClientEntity) GetConfigContent() (*v1.ClientConfig, error) {
	cliCfg, err := utils.LoadClientConfigNormal(c.ConfigContent, true)
	if err != nil {
		return nil, err
	}
	return cliCfg, err
}

func (c *ClientEntity) MarshalJSONConfig() ([]byte, error) {
	cliCfg, err := c.GetConfigContent()
	if err != nil {
		return nil, err
	}
	return json.Marshal(cliCfg)
}

func (c *ClientEntity) ToPB() *pb.Client {
	resp := &pb.Client{
		Id:             &c.ClientID,
		Secret:         &c.ConnectSecret,
		Config:         lo.ToPtr(string(c.ConfigContent)),
		Comment:        &c.Comment,
		ServerId:       &c.ServerID,
		Stopped:        &c.Stopped,
		OriginClientId: &c.OriginClientID,
		FrpsUrl:        &c.FrpsUrl,
		Ephemeral:      &c.Ephemeral,
	}
	if c.LastSeenAt != nil {
		resp.LastSeenAt = lo.ToPtr(c.LastSeenAt.UnixMilli())
	}

	return resp
}
