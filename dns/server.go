package dns

import (
	"denis/config"
	"denis/util"
	"fmt"
	"log"
	"net"
	"strconv"
)

func StartDNS(dnsConfig *config.DNSConfig, records *[]config.Record) error {
	socketAddr := "127.0.0.1:" + strconv.Itoa(dnsConfig.Port)
	udpAddress, err := net.ResolveUDPAddr("udp", socketAddr)
	if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return err
	}

	// UDP messages 512 octets (bytes) or less

	go func() {
		defer connection.Close()
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

			header := parseHeader(message[:offset])

			name := parseName(message, &offset)

			fmt.Println("Resolved Domain:", name.Domain)

			// qType := binary.BigEndian.Uint16(message[offset : offset+2])
			offset += 2

			// qClass := binary.BigEndian.Uint16(message[offset : offset+2])
			offset += 2

			record, ok := config.FindRecordByName(*records, name.Domain)
			if !ok {
				forwardQuery(dnsConfig.Upstream, message, connection, clientAddr)
				continue
			}

			_, port, variant, err := util.ParseAddress(record.Value)

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
			sendAnswer(connection,
				clientAddr,
				&header,
				message,
				offset,
				name.Raw,
				record)
		}
	}()

	return nil
}
