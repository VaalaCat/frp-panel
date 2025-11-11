package wg

import (
	"errors"
	"fmt"
	"net"
	"strings"

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
func parseAndValidatePeerConfigs(peerConfigs []*defs.WireGuardPeerConfig) ([]*defs.WireGuardPeerConfig, error) {
	if len(peerConfigs) == 0 {
		return []*defs.WireGuardPeerConfig{}, nil
	}
	wgTypedPeers := make([]*defs.WireGuardPeerConfig, 0, len(peerConfigs))
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
func parseAndValidatePeerConfig(pCfg *defs.WireGuardPeerConfig) (*defs.WireGuardPeerConfig, error) {

	for _, cidrStr := range pCfg.GetAllowedIps() {
		trimmedCidr := strings.TrimSpace(cidrStr)
		if trimmedCidr == "" {
			continue
		}
		// _, ipNet, err := net.ParseCIDR(trimmedCidr)
		_, _, err := net.ParseCIDR(trimmedCidr)
		if err != nil {
			return nil, errors.Join(errors.New("invalid AllowedIP CIDR"), err)
		}
	}

	if pCfg.PersistentKeepalive <= 0 {
		pCfg.PersistentKeepalive = defs.DefaultPersistentKeepalive
	}

	return pCfg, nil
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// generateUAPIConfigString implementation from previous step
func generateUAPIConfigString(cfg *defs.WireGuardConfig,
	wgPrivateKey wgtypes.Key,
	peerConfigs []*defs.WireGuardPeerConfig,
	firstStart bool,
	skipListenPort bool,
) string {
	uapiBuilder := NewUAPIBuilder()
	uapiBuilder.WithPrivateKey(wgPrivateKey)

	// 只在首次启动且未跳过时设置 listen_port,避免在设备运行时更新端口导致死锁
	if firstStart && !skipListenPort {
		uapiBuilder.WithListenPort(int(cfg.ListenPort))
	}

	uapiBuilder.ReplacePeers(!firstStart).
		AddPeers(peerConfigs)

	return uapiBuilder.Build()
}
