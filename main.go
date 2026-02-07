package main

import (
	"denis/config"
	"denis/dns"
)

func main() {
	denisConfig, confpath := config.LoadDNSConfig()
	dnsConfig := &denisConfig.DNS
	records := config.LoadDNSRecords()

	builder := NewServiceBuilder()

	var dnsErr error
	if dnsConfig.Enabled {
		dnsErr = dns.StartDNS(dnsConfig, &records)
	}
	builder.AddDNS(dnsConfig.Enabled, dnsConfig.Port, dnsErr)

	cliWelcome(CLIConfig{
		ConfigPath:  confpath,
		RecordCount: len(records),
		Services:    builder.Build(),
	})

	if !dnsConfig.Enabled && dnsErr != nil {
		return
	}
}
