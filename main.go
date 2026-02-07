package main

import (
	"denis/config"
	"denis/dns"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	// waits for interrupt signal (from ai)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println(red("\nshutting down."))
}
