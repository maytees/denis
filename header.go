package main

import (
	"encoding/binary"
)

type Flags struct {
	QR     bool
	OPCODE uint8
	AA     bool
	TC     bool
	RD     bool
	RA     bool
	Z      uint8 // Should always be 0
	RCODE  uint8
}

type Header struct {
	ID      uint16
	FLAGS   Flags
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func ParseHeader(header []byte) Header {
	return Header{
		ID:      binary.BigEndian.Uint16(header[0:2]),
		FLAGS:   composeFlagFromBytes(header[2:4]),
		QDCount: binary.BigEndian.Uint16(header[4:6]),
		ANCount: binary.BigEndian.Uint16(header[6:8]),
		NSCount: binary.BigEndian.Uint16(header[8:10]),
		ARCount: binary.BigEndian.Uint16(header[10:12]),
	}
}

// Used for responding
func composeFlag(
	QR bool,
	OPCODE uint8,
	AA bool,
	TC bool,
	RD bool,
	RA bool,
	RCODE uint8,
) uint16 {
	var flags uint16 = 0

	// QR, 1 for response, 0 for query
	flags |= boolToUint16(QR) << 15

	// Opcode, 4 bits, 15 - 4 = 11, 0 = standard query
	flags |= uint16(OPCODE) << 11

	// AA, not sure, set to 0
	flags |= boolToUint16(AA) << 10
	// TC, not sure, set to 0
	flags |= boolToUint16(TC) << 9

	// RD, copies bit at pos 8 (RD) from queryFlags
	// flags |= (queryFlags >> 8) & 1 << 8
	flags |= boolToUint16(RD) << 8

	// RA, recursion available
	flags |= boolToUint16(RA) << 7

	// Z, empty, so is this line necessary?
	flags |= 0 << 4

	// RCODE
	flags |= uint16(RCODE)

	return flags
}

// Returns a flag struct from the 16 bit header flag section
// Used for understanding queries
func composeFlagFromBytes(raw []byte) Flags {
	flags := binary.BigEndian.Uint16(raw)

	return Flags{
		QR: ((flags >> 15) & 1) == 1,
		// 0xF means 4
		OPCODE: uint8((flags >> 11) & 0xF),
		AA:     ((flags >> 10) & 1) == 1,
		TC:     ((flags >> 9) & 1) == 1,
		RD:     ((flags >> 8) & 1) == 1,
		RA:     ((flags >> 7) & 1) == 1,
		// 0x7 means 3 bits
		Z:     uint8((flags >> 4) & 0x7),
		RCODE: uint8(flags & 0xF),
	}
}
