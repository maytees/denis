package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
)

func cliWelcome() {
	fmt.Println("Welcome to the DENIS DNS server!")
}

type ServerConfig struct {
	Port int
}

type Header struct {
	ID      uint16
	FLAGS   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func ParseHeader(header []byte) Header {
	return Header{
		ID:      binary.BigEndian.Uint16(header[0:2]),
		FLAGS:   binary.BigEndian.Uint16(header[2:4]),
		QDCount: binary.BigEndian.Uint16(header[4:6]),
		ANCount: binary.BigEndian.Uint16(header[6:8]),
		NSCount: binary.BigEndian.Uint16(header[8:10]),
		ARCount: binary.BigEndian.Uint16(header[10:12]),
	}
}

func main() {
	cliWelcome()

	udpAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:5354")
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
		log.Printf("\n\nFROM \"%v\" (%d bytes)\n%x", clientAddr.String(), input, message)

		offset := 12 // 12 Bytes?
		rawHeader := message[:offset]

		header := ParseHeader(rawHeader)
		fmt.Printf("\nHeader: \n\tID: %v\n\tFlags: %v\n\tQDCOUNT (Question): %v\n\tANCOUNT (Answer): %v\n\tNSCOUNT (Authority): %v\n\tARCOUNT (Additional): %v\n\t\n",
			header.ID,
			header.FLAGS,
			header.QDCount,
			header.ANCount,
			header.NSCount,
			header.ARCount)

		qNameMap := []string{}

		for {
			length := message[offset]
			offset += 1
			qNameMap = append(qNameMap, string(message[offset:(offset+int(length))]))
			offset += int(length)

			if message[offset] == 0 {
				offset += 1 // moves off the 0
				break
			}
		}

		resolvedDomain := strings.Join(qNameMap, ".")
		fmt.Println("Resolved Domain:", resolvedDomain)

		qType := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		qClass := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		fmt.Printf("QType: %x\nQClass: %x\n", qType, qClass)

		fmt.Printf("Sending response back...\n")

		// responseData := make([]byte, 512)
		// x, err := connection.WriteToUDP(responseData, udpAddress)
		// if err != nil {
		// 	log.Print("Error occured when sending response:", err)
		// }

		fmt.Printf("Whole: %x\n", message)
		fmt.Printf("Question?: %x\n", message[offset:])
		fmt.Printf("Answer?: %x\n", message[:offset])
	}
}
