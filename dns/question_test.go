package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQuestionsSingle(t *testing.T) {
	// google.com: 06 "google" 03 "com" 00, QTYPE=1, QCLASS=1
	message := []byte{
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, // "google"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00,       // null terminator
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}
	offset := 0

	questions, err := parseQuestions(message, &offset, 1)

	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, "google.com", questions[0].QName.Domain)
	assert.Equal(t, uint16(1), uint16(questions[0].QType))
	assert.Equal(t, uint16(1), uint16(questions[0].QClass))
	assert.Equal(t, 16, offset)
}

func TestParseQuestionsSingleLabel(t *testing.T) {
	// localhost: 09 "localhost" 00, QTYPE=1, QCLASS=1
	message := []byte{
		0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6C, 0x68, 0x6F, 0x73, 0x74, // "localhost"
		0x00,       // null terminator
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}
	offset := 0

	questions, err := parseQuestions(message, &offset, 1)

	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, "localhost", questions[0].QName.Domain)
	assert.Equal(t, uint16(1), uint16(questions[0].QType))
	assert.Equal(t, uint16(1), uint16(questions[0].QClass))
	assert.Equal(t, 15, offset)
}

func TestParseQuestions(t *testing.T) {
	// google.com: 02 "go" 03 "com" 00, QTYPE=1, QCLASS=1
	message := []byte{
		0x02, 0x67, 0x6f, // "go"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00,       // null terminator
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN

		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, // "google"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00,      // null terminator
		0x00, 0xF, // QTYPE = MX
		0x00, 0x01, // QCLASS = IN
	}
	offset := 0
	questions, err := parseQuestions(message, &offset, 2)

	require.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.Equal(t, "go.com", questions[0].QName.Domain)
	assert.Equal(t, "google.com", questions[1].QName.Domain)
	assert.Equal(t, uint16(1), uint16(questions[0].QType))
	assert.Equal(t, uint16(15), uint16(questions[1].QType))
}

func TestParseQuestionsShortPacket(t *testing.T) {
	message := []byte{0x06, 0x67, 0x6f}
	offset := 0

	_, err := parseQuestions(message, &offset, 1)

	assert.Error(t, err)
}

func TestParseQuestionsEmptyQuestion(t *testing.T) {
	message := []byte{}
	offset := 0

	questions, err := parseQuestions(message, &offset, 0)

	assert.NoError(t, err)
	assert.Empty(t, questions)
}

func TestParseQuestionsValidNameMissingQType(t *testing.T) {
	// google.com: 02 "go" 03 "com" 00, QTYPE=1, QCLASS=1
	message := []byte{
		0x02, 0x67, 0x6f, // "go"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00,       // null terminator
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN

		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, // "google"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00, // null terminator
	}
	offset := 0
	_, err := parseQuestions(message, &offset, 2)

	require.Error(t, err)
}

func TestParseQuestionsOffsetStartsNonZero(t *testing.T) {
	message := []byte{
		// Header (12 bytes)
		0x12, 0x34, // ID
		0x01, 0x00, // Flags
		0x00, 0x01, // QDCOUNT (1 expected)
		0x00, 0x00, // ANCOUNT
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT

		// Question starts at offset 12
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, // "google"
		0x03, 0x63, 0x6f, 0x6d, // "com"
		0x00,       // null
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN (ends @ offset 28)
	}
	offset := 12

	questions, err := parseQuestions(message, &offset, 1)
	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, "google.com", questions[0].QName.Domain)
	assert.Equal(t, uint16(1), uint16(questions[0].QType))
	assert.Equal(t, uint16(1), uint16(questions[0].QClass))
	assert.Equal(t, 28, offset)
}
