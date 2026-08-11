// NOTE: This test file was written by AI (Claude).
// These are functional tests: they run sendAnswer/forwardQuery over real
// UDP sockets on loopback (ephemeral ports, no sudo needed).
package dns

import (
	"denis/config"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func listenLoopbackUDP(t *testing.T) *net.UDPConn {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)
	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestSendAnswerEndToEnd(t *testing.T) {
	serverConn := listenLoopbackUDP(t)
	clientConn := listenLoopbackUDP(t)
	clientAddr := clientConn.LocalAddr().(*net.UDPAddr)

	query := []byte{
		0xAB, 0xCD, // ID
		0x01, 0x00, // Flags: RD
		0x00, 0x01, // QDCOUNT
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // AN/NS/AR

		0x02, 'm', 'y',
		0x06, 'g', 'o', 'o', 'g', 'l', 'e',
		0x00,
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}
	nameRaw := query[12:23]
	questionEnd := len(query) // 27

	header := parseHeader(query[:12])
	record := config.Record{Name: "my.google", Type: "A", Value: "192.168.1.50", TTL: 300}

	sendAnswer(serverConn, clientAddr, &header, query, questionEnd, nameRaw, record)

	buffer := make([]byte, 512)
	require.NoError(t, clientConn.SetReadDeadline(time.Now().Add(2*time.Second)))
	n, _, err := clientConn.ReadFromUDP(buffer)
	require.NoError(t, err, "expected a response on the client socket")
	response := buffer[:n]

	// header 12 + question 15 + name 11 + type 2 + class 2 + ttl 4 + rdlength 2 + rdata 4
	require.Equal(t, 52, n, "response should be trimmed to its actual length")

	// Header: same ID, response flags, exactly one question and one answer
	assert.Equal(t, uint16(0xABCD), binary.BigEndian.Uint16(response[0:2]))
	flags := composeFlagFromBytes(response[2:4])
	assert.True(t, flags.QR, "QR must be set on a response")
	assert.True(t, flags.AA, "DENIS owns the record, AA should be set")
	assert.True(t, flags.RD, "RD should be copied from the query")
	assert.True(t, flags.RA)
	assert.Equal(t, uint8(0), flags.RCODE)
	assert.Equal(t, uint16(1), binary.BigEndian.Uint16(response[4:6]), "QDCOUNT")
	assert.Equal(t, uint16(1), binary.BigEndian.Uint16(response[6:8]), "ANCOUNT")

	// Question section is echoed back unchanged
	assert.Equal(t, query[12:questionEnd], response[12:questionEnd])

	// Answer section
	answer := response[questionEnd:]
	assert.Equal(t, nameRaw, answer[:11], "answer name")
	assert.Equal(t, uint16(1), binary.BigEndian.Uint16(answer[11:13]), "TYPE A")
	assert.Equal(t, uint16(1), binary.BigEndian.Uint16(answer[13:15]), "CLASS IN")
	assert.Equal(t, uint32(300), binary.BigEndian.Uint32(answer[15:19]), "TTL")
	assert.Equal(t, uint16(4), binary.BigEndian.Uint16(answer[19:21]), "RDLENGTH")
	assert.Equal(t, []byte{192, 168, 1, 50}, answer[21:25], "RDATA")
}

// GAP: a record whose value isn't a valid IPv4 (typo in records.toml,
// or an IPv6 address) makes net.ParseIP(...).To4() return nil and
// sendAnswer panics indexing into it. The caller is expected to validate
// first, but sendAnswer itself should not be able to crash the server.
func TestSendAnswerInvalidRecordValue(t *testing.T) {
	serverConn := listenLoopbackUDP(t)
	clientConn := listenLoopbackUDP(t)
	clientAddr := clientConn.LocalAddr().(*net.UDPAddr)

	query := []byte{
		0xAB, 0xCD, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 't', 'e', 's', 't', 0x00,
		0x00, 0x01, 0x00, 0x01,
	}
	nameRaw := query[12:18]
	header := parseHeader(query[:12])
	record := config.Record{Name: "test", Type: "A", Value: "not-an-ip", TTL: 300}

	require.NotPanics(t, func() {
		sendAnswer(serverConn, clientAddr, &header, query, len(query), nameRaw, record)
	})
}

func TestForwardQueryEndToEnd(t *testing.T) {
	upstreamConn := listenLoopbackUDP(t)
	resolverConn := listenLoopbackUDP(t)
	clientConn := listenLoopbackUDP(t)
	clientAddr := clientConn.LocalAddr().(*net.UDPAddr)

	query := []byte{0xBE, 0xEF, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	cannedResponse := []byte{0xBE, 0xEF, 0x81, 0x80, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0xCA, 0xFE}

	// Fake upstream: reply to whatever arrives with the canned response
	received := make(chan []byte, 1)
	go func() {
		buffer := make([]byte, 512)
		upstreamConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, addr, err := upstreamConn.ReadFromUDP(buffer)
		if err != nil {
			close(received)
			return
		}
		received <- append([]byte{}, buffer[:n]...)
		upstreamConn.WriteToUDP(cannedResponse, addr)
	}()

	forwardQuery(upstreamConn.LocalAddr().String(), query, resolverConn, clientAddr)

	// The upstream must receive the query unmodified
	forwarded, ok := <-received
	require.True(t, ok, "upstream never received the forwarded query")
	assert.Equal(t, query, forwarded)

	// The client must receive the upstream's response
	buffer := make([]byte, 1024)
	require.NoError(t, clientConn.SetReadDeadline(time.Now().Add(2*time.Second)))
	n, _, err := clientConn.ReadFromUDP(buffer)
	require.NoError(t, err, "client never received the relayed response")

	// GAP: forwardQuery relays its whole 512-byte buffer instead of the
	// n bytes actually read, so every forwarded response arrives padded
	// with trailing zeros.
	assert.Equal(t, len(cannedResponse), n, "relayed response should be trimmed to its actual length")
	assert.Equal(t, cannedResponse, buffer[:len(cannedResponse)], "relayed bytes should match the upstream response")
}
