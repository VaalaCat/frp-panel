package wg

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type UAPIBuilder struct {
	// Interface-level settings
	interfacePrivateKey *wgtypes.Key
	listenPort          *int
	fwmark              *int
	replacePeers        bool

	// Accumulated peer sections (already formatted key=value lines ending with \n)
	peerSections []string
}

// NewUAPIBuilder creates a new builder instance.
func NewUAPIBuilder() *UAPIBuilder {
	return &UAPIBuilder{peerSections: make([]string, 0, 8)}
}

// WithPrivateKey sets the interface private key.
func (b *UAPIBuilder) WithPrivateKey(key wgtypes.Key) *UAPIBuilder {
	b.interfacePrivateKey = &key
	return b
}

// WithListenPort sets the interface listen port.
func (b *UAPIBuilder) WithListenPort(port int) *UAPIBuilder {
	b.listenPort = &port
	return b
}

// WithFwmark sets the fwmark. Passing 0 indicates removal per UAPI.
func (b *UAPIBuilder) WithFwmark(mark int) *UAPIBuilder {
	b.fwmark = &mark
	return b
}

// ReplacePeers controls whether subsequent peers replace existing ones instead of appending.
func (b *UAPIBuilder) ReplacePeers(replace bool) *UAPIBuilder {
	b.replacePeers = replace
	return b
}

// AddPeerConfig appends a peer configuration section.
func (b *UAPIBuilder) AddPeerConfig(peer wgtypes.PeerConfig) *UAPIBuilder {
	b.peerSections = append(b.peerSections, buildPeerSection(peer, ""))
	return b
}

func (b *UAPIBuilder) AddPeers(peers []wgtypes.PeerConfig) *UAPIBuilder {
	for _, peer := range peers {
		b.AddPeerConfig(peer)
	}
	return b
}

// UpdatePeerConfig appends a peer configuration section with update_only=true.
func (b *UAPIBuilder) UpdatePeerConfig(peer wgtypes.PeerConfig) *UAPIBuilder {
	b.peerSections = append(b.peerSections, buildPeerSection(peer, "update_only=true\n"))
	return b
}

// RemovePeerByKey appends a peer removal section.
func (b *UAPIBuilder) RemovePeerByKey(publicKey wgtypes.Key) *UAPIBuilder {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("public_key=%s\n", hex.EncodeToString(publicKey[:])))
	sb.WriteString("remove=true\n")
	b.peerSections = append(b.peerSections, sb.String())
	return b
}

// RemovePeerByHexPublicKey appends a peer removal section using a hex-encoded public key.
func (b *UAPIBuilder) RemovePeerByHexPublicKey(hexPublicKey string) *UAPIBuilder {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("public_key=%s\n", strings.ToLower(strings.TrimSpace(hexPublicKey))))
	sb.WriteString("remove=true\n")
	b.peerSections = append(b.peerSections, sb.String())
	return b
}

// Build renders the final UAPI configuration string.
// It ensures interface-level keys precede all peer-level keys and ends with a blank line.
func (b *UAPIBuilder) Build() string {
	var sb strings.Builder

	if b.interfacePrivateKey != nil {
		sb.WriteString(fmt.Sprintf("private_key=%s\n", hex.EncodeToString(b.interfacePrivateKey[:])))
	}
	if b.listenPort != nil {
		sb.WriteString(fmt.Sprintf("listen_port=%d\n", *b.listenPort))
	}
	if b.fwmark != nil {
		sb.WriteString(fmt.Sprintf("fwmark=%d\n", *b.fwmark))
	}
	if b.replacePeers {
		sb.WriteString("replace_peers=true\n")
	}

	for _, section := range b.peerSections {
		sb.WriteString(section)
	}

	out := sb.String()
	out = strings.TrimSuffix(out, "\n")
	return out + "\n\n"
}

// buildPeerSection converts a wgtypes.PeerConfig into its UAPI text format.
func buildPeerSection(peer wgtypes.PeerConfig, extraHeader string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("public_key=%s\n", hex.EncodeToString(peer.PublicKey[:])))
	if extraHeader != "" {
		sb.WriteString(extraHeader)
	}

	if peer.PresharedKey != nil && !isZeroKey(*peer.PresharedKey) {
		sb.WriteString(fmt.Sprintf("preshared_key=%s\n", hex.EncodeToString(peer.PresharedKey[:])))
	}

	if peer.Endpoint != nil {
		sb.WriteString(fmt.Sprintf("endpoint=%s\n", normalizeEndpoint(peer.Endpoint)))
	}

	if peer.PersistentKeepaliveInterval != nil && *peer.PersistentKeepaliveInterval > 0 {
		secs := int(peer.PersistentKeepaliveInterval.Seconds())
		sb.WriteString(fmt.Sprintf("persistent_keepalive_interval=%d\n", secs))
	}

	if peer.ReplaceAllowedIPs {
		sb.WriteString("replace_allowed_ips=true\n")
	}

	for _, allowedIP := range peer.AllowedIPs {
		sb.WriteString(fmt.Sprintf("allowed_ip=%s\n", allowedIP.String()))
	}

	return sb.String()
}

func normalizeEndpoint(addr *net.UDPAddr) string {
	// UDPAddr.String formats IPv6 as [ip%zone]:port; keep as-is to match UAPI.
	return addr.String()
}

func isZeroKey(key wgtypes.Key) bool {
	var zero wgtypes.Key
	return key == zero
}
