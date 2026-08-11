// NOTE: This test file was written by AI (Claude).
package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A complete, valid A query for "my.google" (ID 0xBEEF, RD set)
func validQuery() []byte {
	return []byte{
		0xBE, 0xEF, // ID
		0x01, 0x00, // Flags: RD
		0x00, 0x01, // QDCOUNT
		0x00, 0x00, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT

		0x02, 'm', 'y',
		0x06, 'g', 'o', 'o', 'g', 'l', 'e',
		0x00,
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}
}

func TestParseMessageValidQuery(t *testing.T) {
	query := validQuery()

	message, questionEnd, err := parseMessage(query)

	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, uint16(0xBEEF), message.Header.ID)
	assert.True(t, message.Header.FLAGS.RD)
	assert.Equal(t, "my.google", message.Question.QName.Domain)
	assert.Equal(t, TypeA, message.Question.QType)
	assert.Equal(t, ClassIN, message.Question.QClass)
	assert.Equal(t, query, message.Raw)
	assert.Equal(t, len(query), questionEnd, "questionEnd should point past the question section")
}

// GAP: truncating a packet mid-name panics instead of erroring — same
// root cause as TestParseNameMissingTerminator (length byte read without
// a bounds check).
func TestParseMessageTruncatedQuestion(t *testing.T) {
	// Valid header claiming 1 question, but the question is cut off
	query := validQuery()[:15]

	require.NotPanics(t, func() {
		message, _, err := parseMessage(query)
		assert.Error(t, err)
		assert.Nil(t, message)
	})
}

// GAP: a header with QDCOUNT=0 makes questions[0] panic (index out of
// range). Reachable from the network — e.g. some tools send question-less
// probe packets.
func TestParseMessageZeroQuestions(t *testing.T) {
	query := validQuery()
	query[5] = 0x00 // QDCOUNT = 0

	require.NotPanics(t, func() {
		message, _, err := parseMessage(query)
		assert.Error(t, err)
		assert.Nil(t, message)
	})
}

// GAP: a UDP packet shorter than 12 bytes panics (slice out of range on
// message[:12]) instead of erroring. Reachable from the network — a single
// runt packet takes down the read-loop goroutine.
func TestParseMessageRuntPacket(t *testing.T) {
	require.NotPanics(t, func() {
		message, _, err := parseMessage([]byte{0x01, 0x02, 0x03})
		assert.Error(t, err)
		assert.Nil(t, message)
	})
}

func TestMessageString(t *testing.T) {
	message, _, err := parseMessage(validQuery())
	require.NoError(t, err)

	s := message.String()

	assert.Contains(t, s, "my.google")
	assert.Contains(t, s, "48879") // ID 0xBEEF in decimal
}
