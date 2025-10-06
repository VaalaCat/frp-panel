package models

import (
	"errors"
	"net/netip"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"
)

type WireGuard struct {
	gorm.Model
	*WireGuardEntity

	Client  *Client  `json:"client,omitempty" gorm:"foreignKey:ClientID;references:ClientID"`
	Network *Network `json:"network,omitempty" gorm:"foreignKey:NetworkID;references:ID"`

	AdvertisedEndpoints []*Endpoint      `json:"advertised_endpoints,omitempty" gorm:"foreignKey:WireGuardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WireGuardLinks      []*WireGuardLink `json:"wireguard_links,omitempty" gorm:"foreignKey:FromWireGuardID;references:ID"`
}

type WireGuardEntity struct {
	Name     string `gorm:"type:varchar(255);uniqueIndex:idx_client_id_name"`
	UserId   uint32 `gorm:"index"`
	TenantId uint32 `gorm:"index"`

	PrivateKey   string `json:"private_key" gorm:"type:varchar(255)"`
	LocalAddress string `json:"local_address" gorm:"type:varchar(255)"`
	ListenPort   uint32 `json:"listen_port" gorm:"uniqueIndex:idx_client_id_listen_port"`
	InterfaceMtu uint32 `json:"interface_mtu"`

	DnsServers GormArray[string] `json:"dns_servers" gorm:"type:varchar(255)"`
	ClientID   string            `gorm:"type:varchar(64);uniqueIndex:idx_client_id_name;uniqueIndex:idx_client_id_listen_port"`
	NetworkID  uint              `gorm:"index"`
	Tags       GormArray[string] `json:"tags" gorm:"type:varchar(255)"`
}

func (*WireGuard) TableName() string {
	return "wireguards"
}

func (w *WireGuard) GetTags() []string {
	return w.Tags
}

func (w *WireGuard) GetID() uint {
	return uint(w.ID)
}

func ParseIPOrCIDRWithNetip(s string) (netip.Addr, netip.Prefix, error) {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix.Addr(), prefix, nil
	}

	if addr, err := netip.ParseAddr(s); err == nil {
		return addr, netip.Prefix{}, nil
	}

	return netip.Addr{}, netip.Prefix{}, errors.New("invalid ip or cidr")
}

func (w *WireGuard) AsBasePeerConfig() (*pb.WireGuardPeerConfig, error) {
	privKey, err := wgtypes.ParseKey(w.PrivateKey)
	if err != nil {
		return nil, errors.Join(errors.New("parse private key error"), err)
	}
	_, localIPPrefix, err := ParseIPOrCIDRWithNetip(w.LocalAddress)
	if err != nil {
		return nil, errors.Join(errors.New("parse local address error"), err)
	}

	localIPPrefixAllowed := netip.PrefixFrom(localIPPrefix.Addr(), 32)

	resp := &pb.WireGuardPeerConfig{
		Id:                  uint32(w.ID),
		ClientId:            w.ClientID,
		UserId:              w.UserId,
		TenantId:            w.TenantId,
		PublicKey:           privKey.PublicKey().String(),
		AllowedIps:          []string{localIPPrefixAllowed.String()},
		PersistentKeepalive: 20,
		Tags:                w.Tags,
	}
	if len(w.AdvertisedEndpoints) > 0 {
		resp.Endpoint = &pb.Endpoint{
			Host: w.AdvertisedEndpoints[0].Host,
			Port: w.AdvertisedEndpoints[0].Port,
		}
	}

	return resp, nil
}

func (w *WireGuard) FromPB(pb *pb.WireGuardConfig) {
	w.Model = gorm.Model{}
	w.WireGuardEntity = &WireGuardEntity{}

	w.Model.ID = uint(pb.GetId())

	w.Name = pb.GetInterfaceName()
	w.UserId = pb.GetUserId()
	w.TenantId = pb.GetTenantId()
	w.PrivateKey = pb.GetPrivateKey()
	w.LocalAddress = pb.GetLocalAddress()
	w.ListenPort = pb.GetListenPort()
	w.InterfaceMtu = pb.GetInterfaceMtu()
	w.DnsServers = GormArray[string](pb.GetDnsServers())
	w.ClientID = pb.GetClientId()
	w.NetworkID = uint(pb.GetNetworkId())
	w.Tags = GormArray[string](pb.GetTags())

	w.AdvertisedEndpoints = make([]*Endpoint, 0, len(pb.GetAdvertisedEndpoints()))
	for _, e := range pb.GetAdvertisedEndpoints() {
		w.AdvertisedEndpoints = append(w.AdvertisedEndpoints, &Endpoint{
			Model: gorm.Model{ID: uint(e.GetId())},
			EndpointEntity: &EndpointEntity{
				Host:        e.GetHost(),
				Port:        e.GetPort(),
				ClientID:    e.GetClientId(),
				WireGuardID: uint(e.GetWireguardId()),
			},
		})
	}
}

func (w *WireGuard) ToPB() *pb.WireGuardConfig {
	return &pb.WireGuardConfig{
		Id:            uint32(w.ID),
		ClientId:      w.ClientID,
		UserId:        uint32(w.UserId),
		TenantId:      uint32(w.TenantId),
		InterfaceName: w.Name,
		PrivateKey:    w.PrivateKey,
		LocalAddress:  w.LocalAddress,
		ListenPort:    w.ListenPort,
		InterfaceMtu:  w.InterfaceMtu,
		DnsServers:    w.DnsServers,
		NetworkId:     uint32(w.NetworkID),
		Tags:          w.Tags,
		AdvertisedEndpoints: lo.Map(w.AdvertisedEndpoints, func(e *Endpoint, _ int) *pb.Endpoint {
			return &pb.Endpoint{
				Id:          uint32(e.ID),
				Host:        e.Host,
				Port:        e.Port,
				ClientId:    e.ClientID,
				WireguardId: uint32(e.WireGuardID),
			}
		}),
	}
}

type Network struct {
	gorm.Model
	*NetworkEntity

	WireGuard []*WireGuard `json:"wireguard,omitempty" gorm:"foreignKey:NetworkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (n *Network) FromPB(pbData *pb.Network) {
	n.Model = gorm.Model{}
	n.NetworkEntity = &NetworkEntity{}

	n.Model.ID = uint(pbData.GetId())
	n.Name = pbData.GetName()
	n.UserId = pbData.GetUserId()
	n.TenantId = pbData.GetTenantId()
	n.CIDR = pbData.GetCidr()
	n.ACL = JSON[*pb.AclConfig]{Data: pbData.GetAcl()}
}

func (n *Network) ToPB() *pb.Network {
	return &pb.Network{
		Id:       uint32(n.ID),
		UserId:   n.UserId,
		TenantId: n.TenantId,
		Name:     n.Name,
		Cidr:     n.CIDR,
		Acl:      n.ACL.Data,
	}
}

func (*Network) TableName() string {
	return "networks"
}

type NetworkEntity struct {
	Name     string `gorm:"type:varchar(255);index"`
	UserId   uint32 `gorm:"index"`
	TenantId uint32 `gorm:"index"`

	CIDR string              `gorm:"type:varchar(255);index"`
	ACL  JSON[*pb.AclConfig] `gorm:"type:text;index"`
}

type Endpoint struct {
	gorm.Model
	*EndpointEntity

	WireGuard *WireGuard `json:"wireguard,omitempty" gorm:"foreignKey:WireGuardID;references:ID"`
	Client    *Client    `json:"client,omitempty" gorm:"foreignKey:ClientID;references:ClientID"`
}

type EndpointEntity struct {
	Host string `gorm:"uniqueIndex:idx_client_id_host_port"`
	Port uint32 `gorm:"uniqueIndex:idx_client_id_host_port"`

	WireGuardID uint   `gorm:"index"`
	ClientID    string `gorm:"type:varchar(255);uniqueIndex:idx_client_id_host_port"`
}

func (*Endpoint) TableName() string {
	return "endpoints"
}

func (e *Endpoint) ToPB() *pb.Endpoint {
	return &pb.Endpoint{
		Id:          uint32(e.ID),
		Host:        e.Host,
		Port:        e.Port,
		ClientId:    e.ClientID,
		WireguardId: uint32(e.WireGuardID),
	}
}

func (e *Endpoint) FromPB(pbData *pb.Endpoint) {
	e.Model = gorm.Model{}
	e.EndpointEntity = &EndpointEntity{}

	e.Model.ID = uint(pbData.GetId())
	e.Host = pbData.GetHost()
	e.Port = pbData.GetPort()
	e.ClientID = pbData.GetClientId()
	e.WireGuardID = uint(pbData.GetWireguardId())
}
