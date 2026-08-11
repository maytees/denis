// NOTE: This test file was written by AI (Claude).
// Extends the hand-written cases in question_test.go.
package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQuestionsUppercaseLowered(t *testing.T) {
	message := []byte{
		0x02, 'M', 'Y',
		0x06, 'G', 'o', 'O', 'g', 'L', 'e',
		0x00,
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}
	offset := 0

	questions, err := parseQuestions(message, &offset, 1)

	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, "my.google", questions[0].QName.Domain)
}

func TestParseQuestionsAnyTypeAndClass(t *testing.T) {
	message := []byte{
		0x04, 't', 'e', 's', 't',
		0x00,
		0x00, 0xFF, // QTYPE = * (ANY)
		0x00, 0xFF, // QCLASS = * (ANY)
	}
	offset := 0

	questions, err := parseQuestions(message, &offset, 1)

	require.NoError(t, err)
	require.Len(t, questions, 1)
	assert.Equal(t, QTypeAll, questions[0].QType)
	assert.Equal(t, QClassAny, questions[0].QClass)
}

func TestParseQuestionsMissingQClass(t *testing.T) {
	// Name and QTYPE present, message ends before QCLASS
	message := []byte{
		0x04, 't', 'e', 's', 't',
		0x00,
		0x00, 0x01, // QTYPE = A
	}
	offset := 0

	_, err := parseQuestions(message, &offset, 1)

	assert.Error(t, err)
}

func TestParseQuestionsKeepsRawNameBytes(t *testing.T) {
	nameBytes := []byte{
		0x02, 'm', 'y',
		0x06, 'g', 'o', 'o', 'g', 'l', 'e',
		0x00,
	}
	message := append(append([]byte{}, nameBytes...), 0x00, 0x01, 0x00, 0x01)
	offset := 0

	questions, err := parseQuestions(message, &offset, 1)

	require.NoError(t, err)
	require.Len(t, questions, 1)
	// Raw wire-format bytes are needed later to build the answer
	assert.Equal(t, nameBytes, questions[0].QName.Raw)
}
