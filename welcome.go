// The welcome cli *BUILDER* development was AIDED by AI.
//
// This CLI welcome uses a builder pattern to display
// what DENIS is running, and how. It's modular this way
// because the plan is to add more parts to DENIS. Those
// being -> (api, web interface, reverse proxy)
package main

import (
	"fmt"
	"math/rand"
	"os"

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

var welcomeMessages = []string{
	"Denis knows where everything lives.",
	"Your network, named your way.",
	"Local names for local things.",
	"No more memorizing ports.",
	"Because 192.168.1.47:8096 is not a name.",
	"Your domains. Your machine. Your rules.",
	"DNS for the rest of us.",
	"Giving your services names they deserve.",
	"Home is where the hostname is.",
	"Making local feel less lonely.",
	"Names you'll actually remember.",
	"The phonebook for your homelab.",
	"localhost, but friendlier.",
	"Finally, URLs that make sense.",
	"Your services called by name, not number.",
	"Where every service has a home.",
	"Turning IPs into identities.",
	"Your homelab's receptionist.",
	"No more sticky notes with port numbers.",
	"DNS without the drama.",
	"Local routing made simple.",
	"Naming things so you don't have to remember numbers.",
	"Your network's nickname generator.",
	"Because typing ports is so 2005.",
	"The friendly neighborhood name resolver.",
	"Making self-hosting less hostile.",
	"Where domains come home.",
	"Your services deserve better than localhost:3847.",
	"Putting the home in homelab.",
	"One binary to name them all.",
	"Ctrl+C your port memorization.",
	"Local DNS that actually makes sense.",
	"Your network, finally organized.",
	"Stop bookmarking localhost.",
	"Because you have better things to remember.",
	"Giving localhost a personality.",
	"The name game, but for servers.",
	"Your services, your names, your way.",
	"Making 127.0.0.1 feel like home.",
	"Where ports go to retire.",
	"DNS for humans.",
	"Your personal internet phonebook.",
	"No cloud required.",
	"Naming things is hard. Denis makes it easy.",
	"Your homelab's identity crisis, solved.",
	"Local domains without the hassle.",
	"Because every service deserves a real name.",
	"The missing piece of your homelab.",
	"Your network's naming ceremony.",
	"Stop typing, start naming.",
	"Where services get their identity.",
	"DNS that lives on your machine.",
	"Making self-hosting make sense.",
	"Your router's new best friend.",
	"Local names, global convenience.",
	"No more port roulette.",
	"Your services, now with names.",
	"The easy button for local DNS.",
	"Bringing order to localhost chaos.",
	"Because numbers are for computers.",
	"Your homelab's directory service.",
	"Local first, always.",
	"Making your browser happy.",
	"Where localhost gets a promotion.",
	"Your network's yellow pages.",
	"Simplifying the self-hosted life.",
	"Names over numbers.",
	"Your services, properly addressed.",
	"The DNS server that gets you.",
	"Local hosting, global naming.",
	"Your machine, your rules.",
	"Making ports obsolete since 2025.",
	"The homelab whisperer.",
	"Your network's name tag.",
	"Because you shouldn't need a spreadsheet for ports.",
	"Local DNS done right.",
	"Your services, finally findable.",
	"The anti-port movement starts here.",
	"Giving your homelab a vocabulary.",
	"Where every port gets a name.",
	"Your network's new brain.",
	"Making localhost civilized.",
	"The name layer your homelab needed.",
	"Your services, now addressable.",
	"No more guessing games.",
	"DNS for the self-hosted soul.",
	"Your homelab's contact list.",
	"Where IPs become words.",
	"Local names, zero latency.",
	"Your personal naming authority.",
	"Because jellyfin.home beats 192.168.1.23:8096.",
	"The homelab quality of life upgrade.",
	"Your network, actually navigable.",
	"Making local feel premium.",
	"Where services get surnames.",
	"Your DNS, your data.",
	"No corporate middleman.",
	"The localhost upgrade you deserved.",
	"Your network's address book.",
	"Making self-hosting accessible.",
	"Where ports become paths.",
	"Your homelab, humanized.",
	"Local routing without the headache.",
	"The DNS server that stays home.",
	"Your services, properly introduced.",
	"Because good names make good neighbors.",
	"Your network's naming convention.",
	"Making machines memorable.",
	"Where localhost learns names.",
	"Your homelab's GPS.",
	"Local DNS, global peace of mind.",
	"The self-hoster's companion.",
	"Your services, now with directions.",
	"Because nobody memorizes phone numbers anymore either.",
	"Your network's babel fish.",
	"Making ports passé.",
	"Where every service gets a door.",
	"Your homelab's front desk.",
	"Local names for local heroes.",
	"The naming convention you control.",
	"Your services, human readable.",
	"Because life's too short for port numbers.",
	"Your network's name registry.",
	"Making localhost legendary.",
}

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
		Address: fmt.Sprintf("127.0.0.1:%d", port),
		Error:   err,
	})

	return b
}

func (b *ServiceBuilder) Build() []ServiceStatus {
	return b.services
}

func cliWelcome(config CLIConfig) {
	defer color.Unset()

	color.RGB(rand.Intn(255), rand.Intn(255), rand.Intn(255)).Printf("DENIS - %s\n", welcomeMessages[rand.Intn(len(welcomeMessages))])

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
	successCount := 0
	for _, service := range config.Services {
		if service.Enabled {
			enabledCount++

			if service.Error == nil {
				successCount++
			}

			printServiceStatus(service)
			continue
		}

		printServiceDisabled(service)
	}

	// final status
	fmt.Println()
	if enabledCount == 0 {
		fmt.Println(red("aborting."))
		fmt.Println(gray("no services enabled\n"))

		os.Exit(1)
	} else if successCount == 0 {
		fmt.Println(red("aborting."))
		fmt.Println(gray("all services failed to start\n"))

		os.Exit(1)
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
