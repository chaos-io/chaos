package net

import (
	"errors"
	"net"
	"os"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetHostIP get host ipv4 address
func GetHostIP() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ips, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}

	// Filter and print the IPv4 addresses
	for _, ip := range ips {
		ipAddr := net.ParseIP(ip)
		if ipAddr.To4() != nil {
			return ipAddr.String(), nil
		}
	}

	return "", errors.New("can't get host ip")
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	// Filter and print the IPv4 addresses
	for _, addr := range addrs {
		ipAddr := addr.(*net.IPNet)
		if ipAddr.IP.To4() != nil && !ipAddr.IP.IsLoopback() {
			return ipAddr.IP.String(), nil
		}
	}

	return "", errors.New("can't get local ip")
}

func GetEnvIP() (string, error) {
	ip := os.Getenv("HOST_IP")
	if ip == "" {
		return "", errors.New("can't get env ip")
	}

	return ip, nil
}
