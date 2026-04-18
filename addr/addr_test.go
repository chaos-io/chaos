package addr

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeIPsRemovesUnusableAndDuplicates(t *testing.T) {
	ips := normalizeIPs([]net.Addr{
		ipNet("10.0.0.8/24"),
		&net.IPAddr{IP: net.ParseIP("10.0.0.8")},
		ipNet("127.0.0.1/8"),
		ipNet("0.0.0.0/0"),
		ipNet("224.0.0.1/24"),
		ipNet("2001:db8::1/64"),
	})

	require.Len(t, ips, 3)
	assert.Equal(t, "10.0.0.8", ips[0].String())
	assert.Equal(t, "127.0.0.1", ips[1].String())
	assert.Equal(t, "2001:db8::1", ips[2].String())
}

func TestSelectPrimarySelection(t *testing.T) {
	testCases := []struct {
		name string
		ips  []net.IP
		want string
	}{
		{
			name: "prefers private ipv4",
			ips: []net.IP{
				net.ParseIP("127.0.0.1"),
				net.ParseIP("2001:db8::1"),
				net.ParseIP("8.8.8.8"),
				net.ParseIP("10.0.0.8"),
			},
			want: "10.0.0.8",
		},
		{
			name: "falls back to public ipv4",
			ips: []net.IP{
				net.ParseIP("127.0.0.1"),
				net.ParseIP("2001:db8::1"),
				net.ParseIP("8.8.8.8"),
			},
			want: "8.8.8.8",
		},
		{
			name: "falls back to ipv6 before loopback",
			ips: []net.IP{
				net.ParseIP("127.0.0.1"),
				net.ParseIP("2001:db8::1"),
			},
			want: "2001:db8::1",
		},
		{
			name: "falls back to ipv4 loopback",
			ips: []net.IP{
				net.ParseIP("::1"),
				net.ParseIP("127.0.0.1"),
			},
			want: "127.0.0.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := selectPrimary(tc.ips)
			require.NoError(t, err)
			assert.Equal(t, tc.want, ip.String())
		})
	}
}

func TestSelectPrimaryReturnsErrNoAddressFound(t *testing.T) {
	ip, err := selectPrimary([]net.IP{
		net.ParseIP("0.0.0.0"),
		net.ParseIP("ff02::1"),
	})

	assert.Nil(t, ip)
	assert.ErrorIs(t, err, ErrNoAddressFound)
}

func TestSelectLoopbackPrefersIPv4(t *testing.T) {
	ip, err := selectLoopback([]net.IP{
		net.ParseIP("10.0.0.8"),
		net.ParseIP("::1"),
		net.ParseIP("127.0.0.1"),
	})

	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", ip.String())
}

func TestNormalizeHost(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: " ", want: ""},
		{name: "plain host", input: "localhost", want: "localhost"},
		{name: "ipv4 with port", input: "127.0.0.1:8080", want: "127.0.0.1"},
		{name: "ipv6 with port", input: "[::1]:8080", want: "::1"},
		{name: "wrapped ipv6", input: "[::1]", want: "::1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, normalizeHost(tc.input))
		})
	}
}

func TestFromEnv(t *testing.T) {
	t.Run("rejects empty key", func(t *testing.T) {
		ip, err := FromEnv("")
		assert.Nil(t, ip)
		assert.ErrorIs(t, err, ErrInvalidEnvKey)
	})

	t.Run("rejects unset env", func(t *testing.T) {
		ip, err := FromEnv("HOST_IP")
		assert.Nil(t, ip)
		assert.ErrorIs(t, err, ErrEnvAddressNotSet)
	})

	t.Run("rejects invalid ip", func(t *testing.T) {
		t.Setenv("HOST_IP", "not-an-ip")
		ip, err := FromEnv("HOST_IP")
		assert.Nil(t, ip)
		assert.ErrorIs(t, err, ErrInvalidAddress)
	})

	t.Run("returns parsed ip", func(t *testing.T) {
		t.Setenv("HOST_IP", "10.0.0.8")
		ip, err := FromEnv("HOST_IP")
		require.NoError(t, err)
		assert.Equal(t, "10.0.0.8", ip.String())
	})
}

func TestIsLocal(t *testing.T) {
	assert.True(t, IsLocal("localhost"))
	assert.True(t, IsLocal("localhost:8080"))
	assert.False(t, IsLocal("example.com"))
	assert.False(t, IsLocal("8.8.8.8:53"))
}

func TestListReturnsUsableNormalizedIPs(t *testing.T) {
	ips, err := List()
	require.NoError(t, err)
	require.NotEmpty(t, ips)

	seen := map[string]struct{}{}
	for _, ip := range ips {
		require.True(t, isUsableIP(ip))
		key := ip.String()
		_, exists := seen[key]
		assert.False(t, exists)
		seen[key] = struct{}{}
	}
}

func TestPrimaryReturnsAddressFromList(t *testing.T) {
	ips, err := List()
	require.NoError(t, err)
	require.NotEmpty(t, ips)

	ip, err := Primary()
	require.NoError(t, err)
	assert.True(t, containsIP(ips, ip))
}

func TestLoopbackMatchesAvailableState(t *testing.T) {
	ips, err := List()
	require.NoError(t, err)

	hasLoopback := false
	for _, ip := range ips {
		if ip.IsLoopback() {
			hasLoopback = true
			break
		}
	}

	ip, err := Loopback()
	if !hasLoopback {
		assert.Nil(t, ip)
		assert.ErrorIs(t, err, ErrNoAddressFound)
		return
	}

	require.NoError(t, err)
	assert.True(t, ip.IsLoopback())
}

func TestHostIPv4ReturnsIPv4OrNotFound(t *testing.T) {
	ip, err := HostIPv4(context.Background())
	if err == nil {
		require.NotNil(t, ip)
		assert.NotNil(t, ip.To4())
		return
	}

	assert.True(t, errors.Is(err, ErrNoAddressFound) || strings.Contains(err.Error(), "lookup host") || strings.Contains(err.Error(), "hostname"))
}

func TestFreeTCPPortReturnsPositivePortOrSkipsWhenForbidden(t *testing.T) {
	port, err := FreeTCPPort()
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") || strings.Contains(err.Error(), "permission denied") {
			t.Skipf("environment forbids tcp bind: %v", err)
		}
	}

	require.NoError(t, err)
	assert.Positive(t, port)
}

func ipNet(cidr string) *net.IPNet {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	network.IP = ip
	return network
}
