// NOTE: This test file was written by AI (Claude).
package dns

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaderStandardQuery(t *testing.T) {
	// Typical dig query header: ID 0x1234, RD set, 1 question
	header := parseHeader([]byte{
		0x12, 0x34, // ID
		0x01, 0x00, // Flags: RD only
		0x00, 0x01, // QDCOUNT
		0x00, 0x00, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT
	})

	assert.Equal(t, uint16(0x1234), header.ID)
	assert.False(t, header.FLAGS.QR)
	assert.True(t, header.FLAGS.RD)
	assert.False(t, header.FLAGS.AA)
	assert.Equal(t, uint16(1), header.QDCount)
	assert.Equal(t, uint16(0), header.ANCount)
	assert.Equal(t, uint16(0), header.NSCount)
	assert.Equal(t, uint16(0), header.ARCount)
}

func TestParseHeaderResponseCounts(t *testing.T) {
	header := parseHeader([]byte{
		0xAB, 0xCD, // ID
		0x81, 0x83, // Flags: QR RD RA, RCODE 3 (NXDOMAIN)
		0x00, 0x01, // QDCOUNT
		0x00, 0x02, // ANCOUNT
		0x00, 0x01, // NSCOUNT
		0x00, 0x03, // ARCOUNT
	})

	assert.Equal(t, uint16(0xABCD), header.ID)
	assert.True(t, header.FLAGS.QR)
	assert.True(t, header.FLAGS.RD)
	assert.True(t, header.FLAGS.RA)
	assert.Equal(t, uint8(3), header.FLAGS.RCODE)
	assert.Equal(t, uint16(1), header.QDCount)
	assert.Equal(t, uint16(2), header.ANCount)
	assert.Equal(t, uint16(1), header.NSCount)
	assert.Equal(t, uint16(3), header.ARCount)
}

func TestComposeFlagFromBytesOpcodeAndZ(t *testing.T) {
	// 0x7E70: opcode 15 (all 4 bits), AA, TC, Z=7, RCODE 0
	flags := composeFlagFromBytes([]byte{0x7E, 0x70})

	assert.Equal(t, uint8(0xF), flags.OPCODE)
	assert.True(t, flags.AA)
	assert.True(t, flags.TC)
	assert.False(t, flags.RD)
	assert.Equal(t, uint8(0x7), flags.Z)
	assert.Equal(t, uint8(0), flags.RCODE)
}

// composeFlag and composeFlagFromBytes should be inverses of each other
func TestComposeFlagRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		qr     bool
		opcode uint8
		aa     bool
		tc     bool
		rd     bool
		ra     bool
		rcode  uint8
	}{
		{"standard query", false, 0, false, false, true, false, 0},
		{"authoritative answer", true, 0, true, false, true, true, 0},
		{"nxdomain response", true, 0, false, false, true, true, 3},
		{"inverse query truncated", false, 1, false, true, false, false, 0},
		{"server failure", true, 2, false, false, false, true, 2},
		{"everything set", true, 15, true, true, true, true, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := composeFlag(tt.qr, tt.opcode, tt.aa, tt.tc, tt.rd, tt.ra, tt.rcode)

			raw := make([]byte, 2)
			binary.BigEndian.PutUint16(raw, composed)
			parsed := composeFlagFromBytes(raw)

			assert.Equal(t, tt.qr, parsed.QR)
			assert.Equal(t, tt.opcode, parsed.OPCODE)
			assert.Equal(t, tt.aa, parsed.AA)
			assert.Equal(t, tt.tc, parsed.TC)
			assert.Equal(t, tt.rd, parsed.RD)
			assert.Equal(t, tt.ra, parsed.RA)
			assert.Equal(t, tt.rcode, parsed.RCODE)
			assert.Equal(t, uint8(0), parsed.Z)
		})
	}
}

// GAP: parseHeader panics on a short slice instead of erroring. Reachable
// from the network: any UDP packet under 12 bytes takes down the read loop.
func TestParseHeaderShortInput(t *testing.T) {
	require.NotPanics(t, func() {
		parseHeader([]byte{0x12, 0x34})
	})
}
