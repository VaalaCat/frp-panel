package wg

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

type ACLEntity interface {
	GetTags() []string
	GetID() uint
}

type ACL struct {
	*pb.AclConfig
	Links map[uint]uint
	mu    sync.RWMutex
}

func NewACL() *ACL {
	return &ACL{
		AclConfig: &pb.AclConfig{},
	}
}

func (a *ACL) AddRule(sourceTag, destTag string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Acls = append(a.Acls, &pb.AclRuleConfig{Src: []string{sourceTag}, Dst: []string{destTag}})
}

func (a *ACL) CanConnect(src, dst ACLEntity) bool {
	if a == nil || lo.IsNil(a) {
		return true
	}
	if lo.IsNil(a.AclConfig) || len(a.AclConfig.Acls) == 0 {
		return true
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, sTag := range src.GetTags() {
		for _, dTag := range dst.GetTags() {
			if a.matchRule(sTag, dTag) {
				return true
			}
		}
	}
	return false
}

func (a *ACL) matchRule(sourceTag, destTag string) bool {
	for _, r := range a.Acls {
		if lo.Contains(r.Src, sourceTag) && lo.Contains(r.Dst, destTag) {
			switch r.Action {
			case "accept":
				return true
			case "deny":
				return false
			default:
				return true
			}
		}
	}
	return false
}

func (a *ACL) LoadFromJSON(data []byte) error {
	var cfg pb.AclConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("invalid ACL JSON: %w", err)
	}

	for _, rule := range cfg.Acls {
		if strings.ToLower(rule.Action) != "accept" {
			continue
		}
		for _, s := range rule.Src {
			for _, d := range rule.Dst {
				a.AddRule(s, d)
			}
		}
	}
	return nil
}

func (a *ACL) LoadFromPB(cfg *pb.AclConfig) *ACL {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.AclConfig = cfg
	return a
}

type DefaultEntity struct {
	ID string
}
