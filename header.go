package main

import (
	"encoding/binary"
)

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
