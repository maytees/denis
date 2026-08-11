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
	socketAddr := dnsConfig.Host + ":" + strconv.Itoa(dnsConfig.Port)
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

			rawMessage := buffer[:input]
			fmt.Printf("\nFROM \"%v\" (%d bytes)\n", clientAddr.String(), input)

			message, questionEnd, err := parseMessage(rawMessage)
			if err != nil {
				fmt.Println("Failed to parse message: ", err)
				continue
			}

			fmt.Printf("Got message: %v\n", message)

			record, ok := config.FindRecordByName(*records, message.Question.QName.Domain)
			if !ok {
				forwardQuery(dnsConfig.Upstream, rawMessage, connection, clientAddr)
				continue
			}

			_, port, variant, err := util.ParseAddress(record.Value)
			if err != nil {
				fmt.Println("Could not parse address:", err)
				continue
			}
			if variant != "v4" {
				fmt.Println("Record should be an IPv4 address!")
				continue
			}
			if port != "" {
				fmt.Println("Record should not have port number!")
				continue
			}

			sendAnswer(connection, clientAddr, &message.Header, rawMessage,
				questionEnd, message.Question.QName.Raw, record)
		}
	}()

	return nil
}
