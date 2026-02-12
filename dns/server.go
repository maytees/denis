package dns

import (
	"denis/config"
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

			rawMessage := buffer[:input]
			fmt.Printf("\nFROM \"%v\" (%d bytes)\n", clientAddr.String(), input)

			message, err := parseMessage(
				rawMessage,
				records,
				*dnsConfig,
				connection,
				clientAddr,
			)
			if err != nil {
				fmt.Println("Failed to parse message:", err)
			}

			fmt.Printf("Got message: %v", message)

			// Offset sent here plain beacuse it's the end of the question
			// NOTE: If later on authority and additional are implemented before this call
			// Use another var, since offset would be different
			// sendAnswer(connection,
			// 	clientAddr,
			// 	&header,
			// 	message,
			// 	offset,
			// 	name.Raw,
			// 	record)
		}
	}()

	return nil
}

func parseMessage(message []byte, records *[]config.Record, dnsConfig config.DNSConfig, connection *net.UDPConn, clientAddr *net.UDPAddr) (*Message, error) {
	offset := 12
	header := parseHeader(message[:offset])

	_, err := parseQuestions(message, &offset, header.QDCount)
	if err != nil {
		return nil, err
	}

	// NOTE: do smthn i forgot what
	// TODO: handle authority & additional if necessary

	// TODO: Get by record & name
	// record, ok := config.FindRecordByName(*records, name.Domain)
	// if !ok {
	// 	forwardQuery(dnsConfig.Upstream, message, connection, clientAddr)
	// 	return
	// }

	// _, port, variant, err := util.ParseAddress(record.Value)

	// if err != nil {
	// 	fmt.Println("Could not parse address: ", err)
	// 	return
	// }

	// if variant != "v4" {
	// 	fmt.Println("Record should be an IPv4 address!")
	// 	return
	// }

	// if port != "" {
	// 	fmt.Println("Record should not have port number!")
	// 	return
	// }

	// return &Message{
	// 	Header:     header,
	// 	Question:   nil,
	// 	Answer:     nil,
	// 	Authority:  nil,
	// 	Additional: nil,
	// 	Raw:        message,
	// }

	return nil, nil
}
