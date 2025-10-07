package models

import (
	"github.com/VaalaCat/frp-panel/pb"
	"gorm.io/gorm"
)

// WireGuardLink 描述同一 Network 下两个 WireGuard 节点之间的有向链路与其指标。
// 语义：从 FromWireGuardID 指向 ToWireGuardID 的传输路径，
// UpBandwidthMbps 表示从 From -> To 的可用上行带宽；LatencyMs 为单向时延。
// 如需双向链路，请创建两条对向记录。
type WireGuardLink struct {
	gorm.Model
	*WireGuardLinkEntity

	FromWireGuard *WireGuard `json:"from_wireguard,omitempty" gorm:"foreignKey:FromWireGuardID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ToWireGuard   *WireGuard `json:"to_wireguard,omitempty" gorm:"foreignKey:ToWireGuardID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ToEndpoint    *Endpoint  `json:"to_endpoint,omitempty" gorm:"foreignKey:ToEndpointID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WireGuardLinkEntity struct {
	// 多租户
	UserId   uint32 `gorm:"index"`
	TenantId uint32 `gorm:"index"`

	// 归属网络
	NetworkID uint `gorm:"index"`

	// 有向边两端
	FromWireGuardID uint `gorm:"index"`
	ToWireGuardID   uint `gorm:"index"`

	ToEndpointID uint `gorm:"index"`

	// 链路指标
	UpBandwidthMbps   uint32
	DownBandwidthMbps uint32
	LatencyMs         uint32

	// 状态
	Active bool `gorm:"index"`
}

func (*WireGuardLink) TableName() string {
	return "wireguard_links"
}

func (w *WireGuardLink) FromPB(pbData *pb.WireGuardLink) {
	w.Model = gorm.Model{}
	w.WireGuardLinkEntity = &WireGuardLinkEntity{}

	w.Model.ID = uint(pbData.GetId())
	w.FromWireGuardID = uint(pbData.GetFromWireguardId())
	w.ToWireGuardID = uint(pbData.GetToWireguardId())
	w.UpBandwidthMbps = pbData.GetUpBandwidthMbps()
	w.DownBandwidthMbps = pbData.GetDownBandwidthMbps()
	w.LatencyMs = pbData.GetLatencyMs()
	w.Active = pbData.GetActive()
	w.ToEndpointID = uint(pbData.GetToEndpoint().GetId())
}

func (w *WireGuardLink) ToPB() *pb.WireGuardLink {
	return &pb.WireGuardLink{
		Id:                uint32(w.ID),
		FromWireguardId:   uint32(w.FromWireGuardID),
		ToWireguardId:     uint32(w.ToWireGuardID),
		UpBandwidthMbps:   w.UpBandwidthMbps,
		DownBandwidthMbps: w.DownBandwidthMbps,
		LatencyMs:         w.LatencyMs,
		Active:            w.Active,
		ToEndpoint:        w.ToEndpoint.ToPB(),
	}
}
