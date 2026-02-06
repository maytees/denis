// The welcome cli builder development was AIDED by AI.
//
// This CLI welcome uses a builder pattern to display
// what DENIS is running, and how. It's modular this way
// because the plan is to add more parts to DENIS. Those
// being -> (api, web interface, reverse proxy)
package main

import (
	"fmt"

	"github.com/fatih/color"
)

type ServiceStatus struct {
	Name    string
	Enabled bool
	Address string // either ip + port or localhost + port
	Error   error
}

type CLIConfig struct {
	ConfigPath  string
	RecordCount int
	Services    []ServiceStatus
}

// colors for reuse
var (
	gray   = color.New(color.FgHiBlack).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type ServiceBuilder struct {
	services []ServiceStatus
}

func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{
		// At the moment, there is only a DNS service.
		services: make([]ServiceStatus, 0, 1),
	}
}

func (b *ServiceBuilder) AddDNS(enabled bool, port int, err error) *ServiceBuilder {
	b.services = append(b.services, ServiceStatus{
		Name:    "dns",
		Enabled: enabled,
		Address: fmt.Sprintf("0.0.0.1:%d", port),
		Error:   err,
	})

	return b
}

func (b *ServiceBuilder) Build() []ServiceStatus {
	return b.services
}

func cliWelcome(config CLIConfig) {
	defer color.Unset()

	// header
	fmt.Println()
	if config.ConfigPath == "" {
		fmt.Println(red("⚠ could not find config"))
		return
	}

	fmt.Println(gray("loading config from " + config.ConfigPath))

	// dns info

	if config.RecordCount > 0 {
		fmt.Print(gray("found "))
		fmt.Print(color.WhiteString("%d", config.RecordCount))
		fmt.Println(gray(" DNS records"))
	} else {
		fmt.Printf("%s  no DNS records found\n", "⚠")
	}
	fmt.Println()

	// services status
	enabledCount := 0
	for _, service := range config.Services {
		if service.Enabled {
			enabledCount++
			printServiceStatus(service)

			continue
		}

		printServiceDisabled(service)
	}

	// final status
	fmt.Println()
	if enabledCount == 0 {
		fmt.Println(red("aborting."))
		fmt.Println(gray("no services enabled"))
	} else {
		fmt.Println(green("ready"))
	}
	fmt.Println()
}

func printServiceStatus(service ServiceStatus) {
	if service.Error != nil {
		fmt.Printf("%s  %s %s\n", red("✗"), service.Name, red("failed to start"))
		fmt.Println(gray(fmt.Sprintf("    %s", service.Error.Error())))
	} else {
		arrow := green("➜")
		if service.Name == "web ui" {
			arrow = blue("➜")
		}

		address := gray(service.Address)
		if service.Name == "web ui" || service.Name == "api" {
			address = cyan(service.Address)
		}

		fmt.Printf("%s  %s listening on %s\n", arrow, service.Name, address)
	}
}

func printServiceDisabled(service ServiceStatus) {
	fmt.Printf("%s  %s disabled\n", gray("○"), gray(service.Name))
}
