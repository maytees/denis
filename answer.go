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
	// Use dial instead of listenUdp
	connection, err := net.Dial("udp", upstreamAddr)
	if err != nil {
		log.Fatal("Could not dial upstream: ", err)
	}
	defer connection.Close()

	_, err = connection.Write(message)
	if err != nil {
		log.Fatal("Could not forward upstream: ", err)
	}

	response := make([]byte, 512)
	_, err = connection.Read(response)
	if err != nil {
		log.Fatal("Could not read query response: ", err)
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

	// TODO: Hard coded A record
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded IN
	binary.BigEndian.PutUint16(response[offset:], 1)
	offset += 2

	// Hard coded 0 cache TTL
	binary.BigEndian.PutUint32(response[offset:], 0)
	offset += 4 // TODO: 4 or 2?

	// Length
	binary.BigEndian.PutUint16(response[offset:], 4)
	offset += 2

	response[offset] = 127
	response[offset+1] = 0
	response[offset+2] = 0
	response[offset+3] = 1
	offset += 4

	// fmt.Printf("\nResponse: %x\n", response[:offset])

	// The colon before the offset removes all the empty stuff after
	_, err := connection.WriteToUDP(response[:offset], clientAddress)
	if err != nil {
		log.Print("Error occured when sending response:", err)
	}

	// fmt.Println("Sent response back.")
}
