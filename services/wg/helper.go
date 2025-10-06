package wg

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/VaalaCat/frp-panel/defs"
)

func InitAndValidateWGConfig(cfg *defs.WireGuardConfig, logger *logrus.Entry) error {
	logEntry := logger.WithField("op", "initCfg").WithField("iface", cfg.InterfaceName)

	if cfg.PrivateKey == "" || cfg.PrivateKey == defs.PlaceholderPrivateKey {
		return errors.New("private key is required")
	} else {
		_, err := wgtypes.ParseKey(cfg.PrivateKey)
		if err != nil {
			return errors.Join(errors.New("invalid PrivateKey"), err)
		}
		logEntry.Debugf("Using provided PrivateKey. Public key: %s", cfg.GetParsedPublicKey().String())
	}

	if defs.IsZeroKey(cfg.GetParsedPrivKey()) && cfg.PrivateKey != defs.PlaceholderPrivateKey {
		return errors.New("failed to parse and store private key internally")
	}

	if defs.IsZeroKey(cfg.GetParsedPublicKey()) && cfg.PrivateKey != defs.PlaceholderPrivateKey {
		return errors.New("failed to derive and store public key internally")
	}

	if cfg.GetListenPort() == 0 {
		return errors.New("listen port is required")
	}

	if cfg.GetLocalAddress() == "" {
		return errors.New("local address is required")
	}
	// 使用默认MTU
	if cfg.GetInterfaceMtu() == 0 {
		cfg.InterfaceMtu = defs.DefaultDeviceMTU
		logEntry.Debugf("InterfaceMTU using default: %d", cfg.GetInterfaceMtu())
	}

	if _, _, err := net.ParseCIDR(cfg.LocalAddress); err != nil {
		return errors.Join(fmt.Errorf("invalid LocalAddress ('%s')", cfg.LocalAddress), err)
	}
	return nil
}

type WGKeys struct {
	PrivateKeyBase64 string
	PublicKeyBase64  string
	PrivateKey       wgtypes.Key
	PublicKey        wgtypes.Key
}

// GenerateKeys generates a new private key and returns it in base64 format.
// return: private key, public key, private key, public key
func GenerateKeys() WGKeys {
	priv, e := wgtypes.GeneratePrivateKey()
	if e != nil {
		panic(fmt.Errorf("generate private key: %w", e))
	}
	pub := priv.PublicKey()
	return WGKeys{
		PrivateKeyBase64: priv.String(),
		PublicKeyBase64:  pub.String(),
		PrivateKey:       priv,
		PublicKey:        pub,
	}
}

// parseAndValidatePeerConfigs 生成wg UAPI格式的peer配置
func parseAndValidatePeerConfigs(peerConfigs []*defs.WireGuardPeerConfig) ([]wgtypes.PeerConfig, error) {
	if len(peerConfigs) == 0 {
		return []wgtypes.PeerConfig{}, nil
	}
	wgTypedPeers := make([]wgtypes.PeerConfig, 0, len(peerConfigs))
	for i, pCfg := range peerConfigs {
		peerIDForLog := pCfg.ClientId
		if peerIDForLog == "" {
			peerIDForLog = fmt.Sprintf("index %d (PK: %s...)", i, truncate(pCfg.PublicKey, 10))
		} else {
			peerIDForLog = fmt.Sprintf("'%s' (PK: %s...)", pCfg.ClientId, truncate(pCfg.PublicKey, 10))
		}

		typedPeer, err := parseAndValidatePeerConfig(pCfg)
		if err != nil {
			return nil, fmt.Errorf("peer %s: %w", peerIDForLog, err)
		}
		wgTypedPeers = append(wgTypedPeers, typedPeer)
	}
	return wgTypedPeers, nil
}

// parseAndValidatePeerConfig 将frpp使用的PeerConfig转换为wgtypes.PeerConfig，用来给wg设备使用
func parseAndValidatePeerConfig(pCfg *defs.WireGuardPeerConfig) (wgtypes.PeerConfig, error) {
	var typedPeer wgtypes.PeerConfig

	typedPeer.PublicKey = pCfg.GetParsedPublicKey()
	typedPeer.PresharedKey = pCfg.GetParsedPresharedKey()

	typedPeer.AllowedIPs = make([]net.IPNet, 0, len(pCfg.GetAllowedIps()))
	for _, cidrStr := range pCfg.GetAllowedIps() {
		trimmedCidr := strings.TrimSpace(cidrStr)
		if trimmedCidr == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(trimmedCidr)
		if err != nil {
			return wgtypes.PeerConfig{}, errors.Join(errors.New("invalid AllowedIP CIDR"), err)
		}
		typedPeer.AllowedIPs = append(typedPeer.AllowedIPs, *ipNet)
	}

	if pCfg.Endpoint != nil {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", pCfg.Endpoint.Host, pCfg.Endpoint.Port))
		if err != nil {
			return wgtypes.PeerConfig{}, errors.Join(errors.New("invalid endpoint address"), err)
		}
		typedPeer.Endpoint = addr
	}

	if pCfg.PersistentKeepalive <= 0 {
		typedPeer.PersistentKeepaliveInterval = lo.ToPtr(time.Duration(defs.DefaultPersistentKeepalive) * time.Second)
	} else {
		interval := time.Duration(pCfg.PersistentKeepalive) * time.Second
		typedPeer.PersistentKeepaliveInterval = &interval
	}

	typedPeer.ReplaceAllowedIPs = true
	return typedPeer, nil
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// generateUAPIConfigString implementation from previous step
func generateUAPIConfigString(cfg *defs.WireGuardConfig, wgPrivateKey wgtypes.Key, peerConfigs []wgtypes.PeerConfig, firstStart bool) string {
	uapiBuilder := NewUAPIBuilder()
	uapiBuilder.WithPrivateKey(wgPrivateKey).
		WithListenPort(int(cfg.ListenPort)).
		ReplacePeers(!firstStart).
		AddPeers(peerConfigs)

	return uapiBuilder.Build()
}
