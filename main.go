package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func cliWelcome() {
	fmt.Println("Denis DNS server started.")
}

type ServerConfig struct {
	Port int
}

func main() {
	cliWelcome()

	udpAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:53")
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
		// Prints message
		// log.Printf("\n\nFROM \"%v\" (%d bytes)\n%x", clientAddr.String(), input, message)

		log.Printf("\n\nFROM \"%v\" (%d bytes)", clientAddr.String(), input)

		offset := 12 // 12 Bytes?
		rawHeader := message[:offset]

		header := ParseHeader(rawHeader)
		// fmt.Printf("\nHeader: \n\tID: %v\n\tFlags: %v\n\tQDCOUNT (Question): %v\n\tANCOUNT (Answer): %v\n\tNSCOUNT (Authority): %v\n\tARCOUNT (Additional): %v\n\t\n",
		// 	header.ID,
		// 	header.FLAGS,
		// 	header.QDCount,
		// 	header.ANCount,
		// 	header.NSCount,
		// 	header.ARCount)

		qNameMap := []string{}
		nameStart := offset
		nameEnd := -1

		for {
			length := message[offset]
			offset += 1
			qNameMap = append(qNameMap, string(message[offset:(offset+int(length))]))
			offset += int(length)

			if message[offset] == 0 {
				offset += 1 // moves off the 0
				nameEnd = offset
				break
			}
		}

		resolvedDomain := strings.Join(qNameMap, ".")
		nameLabels := message[nameStart:nameEnd]
		// fmt.Printf("Resolved byte name: %x\n", nameLabels)
		fmt.Println("Resolved Domain:", resolvedDomain)

		// qType := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		// qClass := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		// fmt.Printf("QType: %x\nQClass: %x\n", qType, qClass)

		// Offset sent here plain beacuse it's the end of the question
		// TODO: If later on authority and additional are implemented before this call
		// Use another var, since offset would be different
		SendAnswer(connection, clientAddr, &header, message, offset, nameLabels)
	}
}
