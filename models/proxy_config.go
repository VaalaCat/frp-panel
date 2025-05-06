package models

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ProxyConfig struct {
	*gorm.Model
	*ProxyConfigEntity

	WorkerID string `gorm:"type:varchar(255);index"` // 引用的worker
	Worker   Worker
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
	Stopped        bool   `json:"stopped" gorm:"index"`
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

func (p *ProxyConfigEntity) GetTypedProxyConfig() (v1.TypedProxyConfig, error) {
	var cfg v1.TypedProxyConfig
	err := cfg.UnmarshalJSON(p.Content)
	return cfg, err
}

func (p *ProxyConfig) GetTypedProxyConfig() (v1.TypedProxyConfig, error) {
	return p.ProxyConfigEntity.GetTypedProxyConfig()
}

func (p *ProxyConfig) FillClientConfig(cli *ClientEntity) error {
	return p.ProxyConfigEntity.FillClientConfig(cli)
}

func (p *ProxyConfig) FillTypedProxyConfig(cfg v1.TypedProxyConfig) error {
	annotations := cfg.GetBaseConfig().Annotations
	if len(annotations) > 0 {
		if annotations[defs.FrpProxyAnnotationsKey_Ingress] != "" && len(annotations[defs.FrpProxyAnnotationsKey_WorkerId]) > 0 {
			workerId := annotations[defs.FrpProxyAnnotationsKey_WorkerId]
			p.WorkerID = workerId
		}
	}

	return p.ProxyConfigEntity.FillTypedProxyConfig(cfg)
}

func (p *ProxyConfig) ToPB() *pb.ProxyConfig {
	return &pb.ProxyConfig{
		Id:             lo.ToPtr(uint32(p.ID)),
		Name:           lo.ToPtr(p.Name),
		Type:           lo.ToPtr(p.Type),
		Config:         lo.ToPtr(string(p.Content)),
		Stopped:        lo.ToPtr(p.Stopped),
		ServerId:       lo.ToPtr(p.ServerID),
		ClientId:       lo.ToPtr(p.ClientID),
		OriginClientId: lo.ToPtr(p.OriginClientID),
	}
}
