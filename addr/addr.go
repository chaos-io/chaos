// Package addr provides infrastructure helpers for discovering and validating
// local machine addresses.
package addr

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	// ErrNoAddressFound reports that no usable address could be discovered.
	ErrNoAddressFound = errors.New("addr: no address found")
	// ErrEnvAddressNotSet reports that the requested environment variable was not set.
	ErrEnvAddressNotSet = errors.New("addr: environment address not set")
	// ErrInvalidAddress reports that a provided address string is not a valid IP.
	ErrInvalidAddress = errors.New("addr: invalid address")
	// ErrInvalidEnvKey reports that the requested environment variable key is empty.
	ErrInvalidEnvKey = errors.New("addr: invalid environment variable key")
)

// List returns all usable interface IP addresses in normalized form.
func List() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("addr: list interface addresses: %w", err)
	}

	return normalizeIPs(addrs), nil
}

// Primary returns the preferred local address for the current machine.
func Primary() (net.IP, error) {
	return listAndSelect(selectPrimary)
}

// HostIPv4 returns the first IPv4 address resolved from the current host name.
func HostIPv4(ctx context.Context) (net.IP, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("addr: hostname: %w", err)
	}

	values, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("addr: lookup host %q: %w", host, err)
	}

	for _, value := range values {
		ip := parseIP(value)
		if ip != nil && ip.To4() != nil {
			return ip, nil
		}
	}

	return nil, ErrNoAddressFound
}

// Loopback returns the first available loopback address, preferring IPv4.
func Loopback() (net.IP, error) {
	return listAndSelect(selectLoopback)
}

// FromEnv returns a validated IP parsed from the provided environment variable.
func FromEnv(key string) (net.IP, error) {
	if strings.TrimSpace(key) == "" {
		return nil, ErrInvalidEnvKey
	}

	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil, ErrEnvAddressNotSet
	}

	ip := parseIP(value)
	if ip == nil {
		return nil, fmt.Errorf("%w: %q", ErrInvalidAddress, value)
	}

	return ip, nil
}

// IsLocal reports whether the given host or host:port resolves to the local machine.
func IsLocal(value string) bool {
	host := normalizeHost(value)
	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := parseIP(host)
	if ip == nil {
		return false
	}

	ips, err := List()
	if err != nil {
		return false
	}

	return containsIP(ips, ip)
}

// FreeTCPPort allocates a free local TCP port and returns the chosen port number.
func FreeTCPPort() (int, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		return 0, fmt.Errorf("addr: listen tcp: %w", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("%w: listener address %T", ErrInvalidAddress, listener.Addr())
	}

	return addr.Port, nil
}

func normalizeHost(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	host, _, err := net.SplitHostPort(value)
	if err == nil {
		return host
	}

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	}

	return value
}

func normalizeIPs(addrs []net.Addr) []net.IP {
	seen := make(map[string]struct{}, len(addrs))
	ips := make([]net.IP, 0, len(addrs))

	for _, addr := range addrs {
		ip := addrIP(addr)
		if !isUsableIP(ip) {
			continue
		}

		key := ip.String()
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		ips = append(ips, ip)
	}

	return ips
}

func listAndSelect(selectFn func([]net.IP) (net.IP, error)) (net.IP, error) {
	ips, err := List()
	if err != nil {
		return nil, err
	}

	return selectFn(ips)
}

func parseIP(value string) net.IP {
	return canonicalIP(net.ParseIP(strings.TrimSpace(value)))
}

func addrIP(addr net.Addr) net.IP {
	switch value := addr.(type) {
	case *net.IPNet:
		return canonicalIP(value.IP)
	case *net.IPAddr:
		return canonicalIP(value.IP)
	default:
		return nil
	}
}

func canonicalIP(ip net.IP) net.IP {
	if ip == nil {
		return nil
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4
	}

	return ip.To16()
}

func isUsableIP(ip net.IP) bool {
	return ip != nil && !ip.IsUnspecified() && !ip.IsMulticast()
}

func primaryRank(ip net.IP) int {
	if ipv4 := ip.To4(); ipv4 != nil {
		switch {
		case ipv4.IsPrivate():
			return 0
		case ipv4.IsLoopback():
			return 3
		default:
			return 1
		}
	}

	if ip.IsLoopback() {
		return 4
	}

	return 2
}

func selectPrimary(ips []net.IP) (net.IP, error) {
	var best net.IP
	bestRank := 99

	for _, ip := range ips {
		ip = canonicalIP(ip)
		if !isUsableIP(ip) {
			continue
		}

		rank := primaryRank(ip)
		if rank < bestRank {
			best = ip
			bestRank = rank
		}
	}

	if best == nil {
		return nil, ErrNoAddressFound
	}

	return best, nil
}

func selectLoopback(ips []net.IP) (net.IP, error) {
	var ipv6Loopback net.IP

	for _, ip := range ips {
		ip = canonicalIP(ip)
		if !isUsableIP(ip) || !ip.IsLoopback() {
			continue
		}

		if ip.To4() != nil {
			return ip, nil
		}

		if ipv6Loopback == nil {
			ipv6Loopback = ip
		}
	}

	if ipv6Loopback == nil {
		return nil, ErrNoAddressFound
	}

	return ipv6Loopback, nil
}

func containsIP(ips []net.IP, target net.IP) bool {
	target = canonicalIP(target)
	if target == nil {
		return false
	}

	for _, ip := range ips {
		if normalized := canonicalIP(ip); normalized != nil && normalized.Equal(target) {
			return true
		}
	}

	return false
}
