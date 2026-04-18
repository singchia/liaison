// Package firewall provides an in-process source-IP allowlist for proxies.
//
// Each proxy ID is associated with an optional list of allowed CIDRs. When an
// incoming connection is accepted, the data plane (HTTP server, Gatekeeper)
// calls Check(proxyID, clientIP) before forwarding traffic. A proxy with no
// registered rule is treated as allow-all; a proxy with an empty rule is
// deny-all. Matching is done in-process on the accepted TCP connection — we
// deliberately avoid iptables/ipset to stay portable and root-free.
package firewall

import (
	"fmt"
	"net"
	"sync"
)

// Manager holds the per-proxy CIDR allowlist. The zero value is unusable;
// construct one with NewManager.
type Manager struct {
	mu    sync.RWMutex
	rules map[int][]*net.IPNet // proxyID -> allowed CIDRs; empty slice = deny all
}

// NewManager returns an empty manager.
func NewManager() *Manager {
	return &Manager{rules: make(map[int][]*net.IPNet)}
}

// Allow sets the allowlist for proxyID, replacing any existing entry. Parses
// and validates every CIDR; returns an error on the first bad entry without
// touching existing state. An empty cidrs slice is accepted and means deny
// all inbound traffic for this proxy.
func (m *Manager) Allow(proxyID int, cidrs []string) error {
	parsed := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, ipnet, err := net.ParseCIDR(c)
		if err != nil {
			return fmt.Errorf("invalid CIDR %q: %w", c, err)
		}
		parsed = append(parsed, ipnet)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules[proxyID] = parsed
	return nil
}

// Revoke removes the allowlist for proxyID; subsequent Check calls will
// return true (allow-all).
func (m *Manager) Revoke(proxyID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rules, proxyID)
}

// Check returns true if clientIP is allowed to reach proxyID. A proxy with
// no registered rule is allow-all; an empty rule is deny-all.
func (m *Manager) Check(proxyID int, clientIP net.IP) bool {
	if clientIP == nil {
		return true
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	rules, ok := m.rules[proxyID]
	if !ok {
		return true
	}
	for _, r := range rules {
		if r.Contains(clientIP) {
			return true
		}
	}
	return false
}

// CheckAddr is a convenience wrapper: extracts the IP from a net.Addr and
// dispatches to Check. Returns true on Addr parse failure to fail open — a
// broken address should not drop traffic for a proxy with no rule.
func (m *Manager) CheckAddr(proxyID int, addr net.Addr) bool {
	if addr == nil {
		return true
	}
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return true
	}
	ip := net.ParseIP(host)
	return m.Check(proxyID, ip)
}
