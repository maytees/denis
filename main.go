package main

import (
	"encoding/binary"
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
		log.Printf("\n\nFROM \"%v\" (%d bytes)\n%x", clientAddr.String(), input, message)

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
		fmt.Printf("Resolved byte name: %x\n", nameLabels)
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

// TODO: Add options
func composeFlag(queryFlags uint16) (answerFlags uint16) {
	var flags uint16 = 0

	// QR, 1 for response, 0 for query
	flags |= 1 << 15

	// Opcode, 4 bits, 15 - 4 = 11, 0 = standard query
	flags |= 0 << 11

	// AA, not sure, set to 0
	flags |= 0 << 10
	// TC, not sure, set to 0
	flags |= 0 << 9

	// RD, copies bit at pos 8 (RD) from queryFlags
	flags |= (queryFlags >> 8) & 1 << 8

	// RA, recursion available
	flags |= 1 << 7

	// Z, empty, so is this line necessary?
	flags |= 0 << 4

	// TODO: have different status', keep as no error for now (0).
	flags |= 0 // no need to do <<0

	return flags
}

func SendAnswer(connection *net.UDPConn,
	clientAddress *net.UDPAddr,
	queryHeader *Header,
	message []byte,
	questionEndOffset int,
	nameLabels []byte,
) {
	response := make([]byte, 512)
	offset := 0

	// Forming header, 12 bytes
	binary.BigEndian.PutUint16(response[offset:], queryHeader.ID)
	offset += 2

	flags := composeFlag(queryHeader.FLAGS)
	binary.BigEndian.PutUint16(response[offset:], flags)
	offset += 2

	// 1 Question
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// 1 Answer
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// No Authority
	binary.BigEndian.PutUint16(response[offset:], 0)
	offset += 2

	// No Additional
	binary.BigEndian.PutUint16(response[offset:], 0)
	offset += 2

	// Question section, just copy from original query
	copy(response[offset:], message[12:questionEndOffset])
	offset += questionEndOffset - 12

	// Answer - TODO: Don't hardcode
	copy(response[offset:], nameLabels)
	offset += len(nameLabels)

	// Hard coded A record
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded IN
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded 0 cache TTL
	binary.BigEndian.PutUint32(response[offset:], 0)
	offset += 4 // TODO: 4 or 2?

	binary.BigEndian.PutUint16(response[offset:], 4)
	offset += 2

	response[offset] = 127
	response[offset+1] = 0
	response[offset+2] = 0
	response[offset+3] = 2
	offset += 4

	fmt.Printf("\nResponse: %x\n", response[:offset])

	// The colon before the offset removes all the empty stuff after
	_, err := connection.WriteToUDP(response[:offset], clientAddress)
	if err != nil {
		log.Print("Error occured when sending response:", err)
	}

	fmt.Println("Sent response back.")
}
