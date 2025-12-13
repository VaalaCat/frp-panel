package wg

import (
	"testing"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
)

func TestMergeConnectablePeersFromAdj_AddsMissingPeerWithEmptyAllowedIPs(t *testing.T) {
	ifce := &defs.WireGuardConfig{WireGuardConfig: &pb.WireGuardConfig{
		Id: 1,
		Adjs: map[uint32]*pb.WireGuardLinks{
			1: {Links: []*pb.WireGuardLink{
				{ToWireguardId: 2, ToEndpoint: &pb.Endpoint{Host: "1.2.3.4", Port: 51820}},
			}},
		},
	}}

	desired := []*defs.WireGuardPeerConfig{
		{WireGuardPeerConfig: &pb.WireGuardPeerConfig{
			Id:        3,
			PublicKey: "pk-3",
			AllowedIps: []string{
				"10.0.0.3/32",
			},
		}},
	}

	known := []*defs.WireGuardPeerConfig{
		{WireGuardPeerConfig: &pb.WireGuardPeerConfig{
			Id:        2,
			PublicKey: "pk-2",
			AllowedIps: []string{
				"10.0.0.2/32",
			},
			PersistentKeepalive: 0, // should be defaulted by parseAndValidatePeerConfig
		}},
	}

	got := mergeConnectablePeersFromAdj(ifce, desired, known)

	var found *defs.WireGuardPeerConfig
	for _, p := range got {
		if p != nil && p.GetId() == 2 {
			found = p
			break
		}
	}
	if found == nil {
		t.Fatalf("expected peer id=2 to be added")
	}
	if len(found.GetAllowedIps()) != 0 {
		t.Fatalf("expected allowed_ips empty, got=%v", found.GetAllowedIps())
	}
	if found.GetEndpoint() == nil || found.GetEndpoint().GetHost() != "1.2.3.4" || found.GetEndpoint().GetPort() != 51820 {
		t.Fatalf("expected endpoint from adj to be applied, got=%v", found.GetEndpoint())
	}
	if found.GetPersistentKeepalive() == 0 {
		t.Fatalf("expected persistent_keepalive defaulted, got=0")
	}
}

func TestMergeConnectablePeersFromAdj_DoesNotOverrideExistingDesiredPeer(t *testing.T) {
	ifce := &defs.WireGuardConfig{WireGuardConfig: &pb.WireGuardConfig{
		Id: 1,
		Adjs: map[uint32]*pb.WireGuardLinks{
			1: {Links: []*pb.WireGuardLink{
				{ToWireguardId: 2, ToEndpoint: &pb.Endpoint{Host: "9.9.9.9", Port: 9999}},
			}},
		},
	}}

	desired := []*defs.WireGuardPeerConfig{
		{WireGuardPeerConfig: &pb.WireGuardPeerConfig{
			Id:        2,
			PublicKey: "pk-2",
			AllowedIps: []string{
				"10.0.0.2/32",
			},
			Endpoint: &pb.Endpoint{Host: "1.1.1.1", Port: 1111},
		}},
	}

	got := mergeConnectablePeersFromAdj(ifce, desired, nil)

	if len(got) != 1 {
		t.Fatalf("expected no extra peers added, got len=%d", len(got))
	}
	if len(got[0].GetAllowedIps()) != 1 || got[0].GetAllowedIps()[0] != "10.0.0.2/32" {
		t.Fatalf("expected desired allowed_ips preserved, got=%v", got[0].GetAllowedIps())
	}
	// endpoint should not be overridden because peer already existed in desired list
	if got[0].GetEndpoint() == nil || got[0].GetEndpoint().GetHost() != "1.1.1.1" {
		t.Fatalf("expected desired endpoint preserved, got=%v", got[0].GetEndpoint())
	}
}

func TestMergeConnectablePeersFromAdj_UsesEndpointWireguardIDWhenPeerIDMissing(t *testing.T) {
	ifce := &defs.WireGuardConfig{WireGuardConfig: &pb.WireGuardConfig{
		Id: 1,
		Adjs: map[uint32]*pb.WireGuardLinks{
			1: {Links: []*pb.WireGuardLink{
				{ToWireguardId: 2, ToEndpoint: &pb.Endpoint{Host: "2.2.2.2", Port: 2222}},
			}},
		},
	}}

	desired := []*defs.WireGuardPeerConfig{}

	known := []*defs.WireGuardPeerConfig{
		{WireGuardPeerConfig: &pb.WireGuardPeerConfig{
			Id:        0, // 模拟：peer.id 未下发
			PublicKey: "pk-2",
			Endpoint:  &pb.Endpoint{WireguardId: 2, Host: "old", Port: 1},
		}},
	}

	got := mergeConnectablePeersFromAdj(ifce, desired, known)
	if len(got) != 1 {
		t.Fatalf("expected 1 peer added, got len=%d", len(got))
	}
	if got[0].GetPublicKey() != "pk-2" {
		t.Fatalf("expected pk-2, got=%s", got[0].GetPublicKey())
	}
	// endpoint should be overridden by adj's to_endpoint
	if got[0].GetEndpoint() == nil || got[0].GetEndpoint().GetHost() != "2.2.2.2" || got[0].GetEndpoint().GetPort() != 2222 {
		t.Fatalf("expected endpoint from adj, got=%v", got[0].GetEndpoint())
	}
	if len(got[0].GetAllowedIps()) != 0 {
		t.Fatalf("expected allowed_ips empty, got=%v", got[0].GetAllowedIps())
	}
}
