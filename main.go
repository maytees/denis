package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// func cliWelcome(configPath string, recordsLength int, dnsAddr string) {
// 	gray := color.New(color.FgHiBlack).SprintFunc()
// 	green := color.New(color.FgGreen).SprintFunc()
// 	// blue := color.New(color.FgBlue).SprintFunc()
// 	// cyan := color.New(color.FgCyan).SprintFunc()
// 	red := color.New(color.FgRed).SprintFunc()
// 	// black := color.New(color.FgBlack).SprintFunc()
// 	yellow := color.New(color.FgYellow).SprintFunc()
// 	defer color.Unset()

// 	// Header
// 	if configPath == "" {
// 		fmt.Println(red("\ncould not find config"))
// 	}

// 	fmt.Println(gray("\nloading config from ", configPath))

// 	if recordsLength > 0 {
// 		fmt.Print(gray("found "))
// 		fmt.Print(color.WhiteString("%d", recordsLength))
// 		fmt.Print(gray(" DNS records"))
// 	} else {
// 		fmt.Printf("%s  no DNS records found\n", yellow("⚠"))
// 	}

// 	fmt.Println()
// 	fmt.Println()

// 	if dnsAddr == "disabled" {
// 		fmt.Printf(red("→  dns disabled\n"))
// 		fmt.Println()

// 		fmt.Println(red("aborting."))
// 	} else {
// 		fmt.Printf("%s  dns listening on %s\n", green("➜"), gray(fmt.Sprintf("%s", dnsAddr)))
// 		fmt.Println()

// 		fmt.Println(green("ready"))
// 	}

// 	fmt.Println()
// }

func main() {
	config, confpath := loadDNSConfig()
	dnsConfig := config.DNS
	records := loadDNSRecords()

	builder := NewServiceBuilder()

	var dnsErr error
	if dnsConfig.Enabled {
		// TODO: startDNS func
	}
	builder.AddDNS(dnsConfig.Enabled, dnsConfig.Port, dnsErr)

	// cache := cache.New(5*time.Minute, 10*time.Minute)

	cliWelcome(CLIConfig{
		ConfigPath:  confpath,
		RecordCount: len(records),
		Services:    builder.Build(),
	})

	// TODO: temporary
	if !dnsConfig.Enabled {
		return
	}

	socketAddr := "127.0.0.1:" + strconv.Itoa(dnsConfig.Port)
	udpAddress, err := net.ResolveUDPAddr("udp", socketAddr)
	if err != nil {
		log.Fatalln(err)
	}

	connection, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()

	// UDP messages 512 octets (bytes) or less
	buffer := make([]byte, 512)
	// UDP messages 512 octets (bytes) or less

	for {
		input, clientAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from packet %v\n", err)
			continue
		}

		message := buffer[:input]

		log.Printf("\n\nFROM \"%v\" (%d bytes)", clientAddr.String(), input)

		offset := 12

		header := ParseHeader(message[:offset])

		name := parseName(message, &offset)

		fmt.Println("Resolved Domain:", name.Domain)

		// qType := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		// qClass := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		record, ok := findRecordByName(records, name.Domain)
		if !ok {
			forwardQuery(dnsConfig.Upstream, message, connection, clientAddr)
			continue
		}

		_, port, variant, err := parseAddress(record.Value)

		if err != nil {
			log.Fatalln("Could not parse address: ", err)
		}

		if variant != "v4" {
			log.Fatalln("Record should be an IPv4 address!")
		}

		if port != "" {
			log.Fatalln("Record should not have port number!")
		}

		// Offset sent here plain beacuse it's the end of the question
		// NOTE: If later on authority and additional are implemented before this call
		// Use another var, since offset would be different
		SendAnswer(connection,
			clientAddr,
			&header,
			message,
			offset,
			name.Raw,
			record)
	}
}
