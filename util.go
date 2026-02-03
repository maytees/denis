package main

import (
	"net"
	"strings"
)

func getVariant(host string) string {
	if ip := net.ParseIP(host); ip != nil {
		if ip.To4() != nil {
			return "v4"
		}

		return "v6"
	}

	return "domain"
}

func parseAddress(addr string) (host string, port string, variant string, err error) {
	host, port, err = net.SplitHostPort(addr)
	if err == nil {
		variant = getVariant(host)
		return host, port, variant, nil
	}

	// Check if valid ip (v6 or v4)
	if ip := net.ParseIP(addr); ip != nil {
		variant = getVariant(addr)

		return addr, "", variant, nil
	}

	// Hostname without a port
	if !strings.Contains(addr, ":") {
		return addr, "", "domain", nil
	}

	return "", "", "", err
}
