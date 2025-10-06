package wg

import (
	"strconv"
	"strings"

	"github.com/VaalaCat/frp-panel/pb"
)

func ParseWGRunningInfo(raw string) (*pb.WGDeviceRuntimeInfo, error) {
	lines := strings.Split(raw, "\n")
	dev := &pb.WGDeviceRuntimeInfo{Peers: make([]*pb.WGPeerRuntimeInfo, 0, 8)}
	var cur *pb.WGPeerRuntimeInfo

	flushPeer := func() {
		if cur != nil {
			dev.Peers = append(dev.Peers, cur)
			cur = nil
		}
	}

	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			break
		}
		eq := strings.IndexByte(ln, '=')
		if eq <= 0 {
			continue
		}
		k := ln[:eq]
		v := ln[eq+1:]
		switch k {
		case "private_key":
			dev.PrivateKey = v
		case "listen_port":
			if p, err := strconv.ParseUint(v, 10, 32); err == nil {
				dev.ListenPort = uint32(p)
			}
		case "protocol_version":
			if pv, err := strconv.ParseUint(v, 10, 32); err == nil {
				dev.ProtocolVersion = uint32(pv)
			}
		case "errno":
			if e, err := strconv.ParseInt(v, 10, 32); err == nil {
				dev.Errno = int32(e)
			}
		case "public_key":
			// 新 peer 开始
			flushPeer()
			cur = &pb.WGPeerRuntimeInfo{PublicKey: v}
		case "preshared_key":
			if cur != nil {
				cur.PresharedKey = v
			}
		case "allowed_ip":
			if cur != nil {
				cur.AllowedIps = append(cur.AllowedIps, v)
			}
		case "endpoint":
			if cur != nil {
				host, port := parseEndpoint(v)
				cur.EndpointHost = host
				cur.EndpointPort = port
			}
		case "tx_bytes":
			if cur != nil {
				if n, err := strconv.ParseUint(v, 10, 64); err == nil {
					cur.TxBytes = n
				}
			}
		case "rx_bytes":
			if cur != nil {
				if n, err := strconv.ParseUint(v, 10, 64); err == nil {
					cur.RxBytes = n
				}
			}
		case "persistent_keepalive_interval":
			if cur != nil {
				if n, err := strconv.ParseInt(v, 10, 32); err == nil {
					cur.PersistentKeepaliveInterval = uint32(n)
				}
			}
		case "last_handshake_time_nsec":
			if cur != nil {
				if n, err := strconv.ParseUint(v, 10, 64); err == nil {
					cur.LastHandshakeTimeNsec = n
				}
			}
		case "last_handshake_time_sec":
			if cur != nil {
				if n, err := strconv.ParseUint(v, 10, 64); err == nil {
					cur.LastHandshakeTimeSec = n
				}
			}
		default:
			if cur != nil {
				if cur.Extra == nil {
					cur.Extra = make(map[string]string)
				}
				cur.Extra[k] = v
			} else {
				if dev.Extra == nil {
					dev.Extra = make(map[string]string)
				}
				dev.Extra[k] = v
			}
		}
	}
	flushPeer()
	return dev, nil
}

func parseEndpoint(v string) (string, uint32) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", 0
	}
	// IPv6 带 [] 的形式
	if strings.HasPrefix(v, "[") {
		rb := strings.IndexByte(v, ']')
		if rb > 1 && rb+1 < len(v) && v[rb+1] == ':' {
			host := v[1:rb]
			portStr := v[rb+2:]
			if p, err := strconv.ParseUint(portStr, 10, 32); err == nil {
				return host, uint32(p)
			}
			return host, 0
		}
		return v, 0
	}
	// IPv4 或域名 host:port
	idx := strings.LastIndexByte(v, ':')
	if idx <= 0 || idx+1 >= len(v) {
		return v, 0
	}
	host := v[:idx]
	portStr := v[idx+1:]
	if p, err := strconv.ParseUint(portStr, 10, 32); err == nil {
		return host, uint32(p)
	}
	return host, 0
}
