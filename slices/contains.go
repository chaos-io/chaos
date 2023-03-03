package slices

import (
	"bytes"
	"net"
)

// ContainsByte checks if byte slice contains given byte
func ContainsByte(haystack []byte, needle byte) bool {
	return bytes.IndexByte(haystack, needle) != -1
}

// ContainsIP checks if net.IP slice contains given net.IP
func ContainsIP(haystack []net.IP, needle net.IP) bool {
	for _, e := range haystack {
		if e.Equal(needle) {
			return true
		}
	}
	return false
}
