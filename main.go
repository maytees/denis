package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func cliWelcome() {
	fmt.Println("Denis DNS server started.")
}

func main() {
	config := loadDNSConfig()
	dnsConfig := config.DNS
	records := loadDNSRecords().Records

	// cache := cache.New(5*time.Minute, 10*time.Minute)

	if !dnsConfig.Enabled {
		log.Fatalln("DNS is disabled")
		return
	}

	cliWelcome()

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
