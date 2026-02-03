package main

import (
	"net"
	"strings"
)

func parseAddress(addr string) (host string, port string, err error) {
	host, port, err = net.SplitHostPort(addr)
	if err == nil {
		return host, port, nil
	}

	// Check if valid ip (v6 or v4)
	if ip := net.ParseIP(addr); ip != nil {
		return addr, "", nil
	}

	// Hostname without a port
	if !strings.Contains(addr, ":") {
		return addr, "", nil
	}

	return "", "", err
}
