package defs

import (
	"bytes"
	"errors"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireGuardConfig struct {
	*pb.WireGuardConfig

	parsedPublicKey wgtypes.Key `json:"-"`
	parsedPrivKey   wgtypes.Key `json:"-"`
}

func (w *WireGuardConfig) GetParsedPublicKey() wgtypes.Key {
	if IsZeroKey(w.parsedPublicKey) {
		w.parsedPublicKey = w.GetParsedPrivKey().PublicKey()
	}

	return w.parsedPublicKey
}

func (w *WireGuardConfig) GetParsedPrivKey() wgtypes.Key {
	if IsZeroKey(w.parsedPrivKey) {
		var err error
		w.parsedPrivKey, err = wgtypes.ParseKey(w.GetPrivateKey())
		if err != nil {
			panic(errors.Join(errors.New("parse private key error"), err))
		}
	}

	return w.parsedPrivKey
}

func (w *WireGuardConfig) GetParsedPeers() []*WireGuardPeerConfig {
	parsedPeers := make([]*WireGuardPeerConfig, 0, len(w.GetPeers()))
	for _, p := range w.GetPeers() {
		parsedPeers = append(parsedPeers, &WireGuardPeerConfig{WireGuardPeerConfig: p})
	}
	return parsedPeers
}

type WireGuardPeerConfig struct {
	*pb.WireGuardPeerConfig

	parsedPublicKey    wgtypes.Key `json:"-"`
	parsedPresharedKey wgtypes.Key `json:"-"`
}

func (w *WireGuardPeerConfig) GetParsedPublicKey() wgtypes.Key {
	if IsZeroKey(w.parsedPublicKey) {
		var err error
		w.parsedPublicKey, err = wgtypes.ParseKey(w.GetPublicKey())
		if err != nil {
			panic(errors.Join(errors.New("parse public key error"), err))
		}
	}

	return w.parsedPublicKey
}

func (w *WireGuardPeerConfig) GetParsedPresharedKey() *wgtypes.Key {
	if w.GetPresharedKey() == "" {
		return nil
	}

	if IsZeroKey(w.parsedPresharedKey) {
		var err error
		w.parsedPresharedKey, err = wgtypes.ParseKey(w.GetPresharedKey())
		if err != nil {
			panic(errors.Join(errors.New("parse preshared key error"), err))
		}
	}
	return lo.ToPtr(w.parsedPresharedKey)
}

func (w *WireGuardPeerConfig) Equal(other *WireGuardPeerConfig) bool {
	endpointEqual := false
	if w.Endpoint != nil && other.Endpoint != nil {
		endpointEqual = (w.Endpoint.Host == other.Endpoint.Host && w.Endpoint.Port == other.Endpoint.Port)
	} else if w.Endpoint == nil && other.Endpoint == nil {
		endpointEqual = true
	}

	oldExtraIps, newExtraIps := lo.Difference(w.GetAllowedIps(), other.GetAllowedIps())
	allowedIpsEqual := len(oldExtraIps) == 0 && len(newExtraIps) == 0

	return w.Id == other.Id &&
		w.ClientId == other.ClientId &&
		w.UserId == other.UserId &&
		w.TenantId == other.TenantId &&
		w.PublicKey == other.PublicKey &&
		w.PresharedKey == other.PresharedKey &&
		w.PersistentKeepalive == other.PersistentKeepalive &&
		endpointEqual &&
		allowedIpsEqual
}

// IsZeroKey 检查一个 wgtypes.Key 是否是空。
func IsZeroKey(key wgtypes.Key) bool {
	var zeroKey wgtypes.Key
	return bytes.Equal(key[:], zeroKey[:])
}

type WireGuardLink struct {
	*pb.WireGuardLink
}

func (w *WireGuardLink) GetReverse() *WireGuardLink {
	return &WireGuardLink{
		WireGuardLink: &pb.WireGuardLink{
			Id:                w.Id,
			FromWireguardId:   w.ToWireguardId,
			ToWireguardId:     w.FromWireguardId,
			UpBandwidthMbps:   w.DownBandwidthMbps,
			DownBandwidthMbps: w.UpBandwidthMbps,
			LatencyMs:         w.LatencyMs,
			Active:            w.Active,
			ToEndpoint:        w.ToEndpoint,
		},
	}
}
