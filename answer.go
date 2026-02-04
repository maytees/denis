package main

import (
	"encoding/binary"
	"log"
	"net"
)

func forwardQuery(upstreamAddr string,
	message []byte,
	resolverConn *net.UDPConn,
	clientAddress *net.UDPAddr,
) {
	_, port, _, err := parseAddress(upstreamAddr)
	if err != nil {
		log.Fatalln("Invalid upstream addr: ", err)
	}

	if port == "" {
		upstreamAddr += ":53"
	}

	// Use dial instead of listenUdp
	connection, err := net.Dial("udp", upstreamAddr)
	if err != nil {
		log.Fatalln("Could not dial upstream: ", err)
	}
	defer connection.Close()

	_, err = connection.Write(message)
	if err != nil {
		log.Fatalln("Could not forward upstream: ", err)
	}

	response := make([]byte, 512)
	_, err = connection.Read(response)
	if err != nil {
		log.Fatalln("Could not read query response: ", err)
	}

	_, err = resolverConn.WriteToUDP(response, clientAddress)
	if err != nil {
		log.Print("Error occursed when sending response: ", err)
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

	// RCODE (Mockapetrics, p.26)
	// TODO: have different status', keep as no error for now (0).
	flags |= 0

	return flags
}

func SendAnswer(connection *net.UDPConn,
	clientAddress *net.UDPAddr,
	queryHeader *Header,
	message []byte,
	questionEndOffset int,
	nameLabels []byte,
	record Record,
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

	copy(response[offset:], nameLabels)
	offset += len(nameLabels)

	// TODO: Hard coded A record
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded IN
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded 0 cache TTL
	binary.BigEndian.PutUint32(response[offset:], record.TTL)
	offset += 4

	// Length
	binary.BigEndian.PutUint16(response[offset:], 4)
	offset += 2

	// ParseIP returns v6 & v4 together, with first 11 being
	// v6 bytes, and the next 4 being v4, hence, To4() gets the last 4 bytes.
	ip := net.ParseIP(record.Value).To4()
	response[offset] = ip[0]
	response[offset+1] = ip[1]
	response[offset+2] = ip[2]
	response[offset+3] = ip[3]

	offset += 4

	// fmt.Printf("\nResponse: %x\n", response[:offset])

	// The colon before the offset removes all the empty stuff after
	_, err := connection.WriteToUDP(response[:offset], clientAddress)
	if err != nil {
		log.Print("Error occured when sending response:", err)
	}

	// fmt.Println("Sent response back.")
}
