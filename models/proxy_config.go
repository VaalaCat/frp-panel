package models

import (
	"fmt"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"gorm.io/gorm"
)

type ProxyConfig struct {
	*gorm.Model
	*ProxyConfigEntity
}

type ProxyConfigEntity struct {
	ServerID       string `json:"server_id" gorm:"index"`
	ClientID       string `json:"client_id" gorm:"index"`
	Name           string `json:"name" gorm:"index"`
	Type           string `json:"type" gorm:"index"`
	UserID         int    `json:"user_id" gorm:"index"`
	TenantID       int    `json:"tenant_id" gorm:"index"`
	OriginClientID string `json:"origin_client_id" gorm:"index"`
	Content        []byte `json:"content"`
}

func (*ProxyConfig) TableName() string {
	return "proxy_config"
}

func (p *ProxyConfigEntity) FillTypedProxyConfig(cfg v1.TypedProxyConfig) error {
	var err error
	p.Name = cfg.GetBaseConfig().Name
	p.Type = cfg.GetBaseConfig().Type
	p.Content, err = cfg.MarshalJSON()
	return err
}

func (p *ProxyConfigEntity) FillClientConfig(cli *ClientEntity) error {
	if cli == nil {
		return fmt.Errorf("invalid client, client is nil")
	}
	p.ServerID = cli.ServerID
	p.ClientID = cli.ClientID
	p.UserID = cli.UserID
	p.TenantID = cli.TenantID
	p.OriginClientID = cli.OriginClientID
	return nil
}
