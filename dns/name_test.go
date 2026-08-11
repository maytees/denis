// NOTE: This test file was written by AI (Claude).
package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNameMultiLabel(t *testing.T) {
	message := []byte{
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
	}
	offset := 0

	name, err := parseName(message, &offset)

	require.NoError(t, err)
	assert.Equal(t, "www.example.com", name.Domain)
	assert.Equal(t, 17, offset)
	// Raw keeps the wire form, including length bytes and terminator
	assert.Equal(t, message, name.Raw)
}

func TestParseNameLowercasesInput(t *testing.T) {
	message := []byte{
		0x06, 'G', 'o', 'O', 'g', 'L', 'e',
		0x03, 'C', 'O', 'M',
		0x00,
	}
	offset := 0

	name, err := parseName(message, &offset)

	require.NoError(t, err)
	assert.Equal(t, "google.com", name.Domain)
}

func TestParseNameRoot(t *testing.T) {
	// A lone null byte is the root domain
	message := []byte{0x00, 0xFF, 0xFF}
	offset := 0

	name, err := parseName(message, &offset)

	require.NoError(t, err)
	assert.Equal(t, "", name.Domain)
	assert.Equal(t, 1, offset)
	assert.Equal(t, []byte{0x00}, name.Raw)
}

func TestParseNameStartsAtOffset(t *testing.T) {
	message := []byte{
		0xDE, 0xAD, // junk before the name
		0x02, 'g', 'o',
		0x00,
	}
	offset := 2

	name, err := parseName(message, &offset)

	require.NoError(t, err)
	assert.Equal(t, "go", name.Domain)
	assert.Equal(t, 6, offset)
}

func TestParseNameTruncatedLabel(t *testing.T) {
	// Length byte promises 5 chars, only 2 present
	message := []byte{0x05, 'a', 'b'}
	offset := 0

	_, err := parseName(message, &offset)

	assert.Error(t, err)
}

// GAP: an empty message panics (index out of range) instead of
// returning an error — the length byte is read without a bounds check.
func TestParseNameEmptyMessage(t *testing.T) {
	message := []byte{}
	offset := 0

	require.NotPanics(t, func() {
		_, err := parseName(message, &offset)
		assert.Error(t, err)
	})
}

// GAP: a name missing its null terminator panics instead of erroring —
// after the last label the loop reads a length byte past the end.
func TestParseNameMissingTerminator(t *testing.T) {
	message := []byte{0x03, 'a', 'b', 'c'}
	offset := 0

	require.NotPanics(t, func() {
		_, err := parseName(message, &offset)
		assert.Error(t, err)
	})
}

// Documents a known limitation: DNS message compression (RFC 1035 §4.1.4,
// pointer bytes 0xC0..) is not supported. A pointer is misread as a huge
// label length and rejected. This matters once upstream responses (which
// almost always use compression) get parsed.
func TestParseNameCompressionPointerUnsupported(t *testing.T) {
	// 0xC0 0x0C = pointer to offset 12
	message := []byte{0xC0, 0x0C}
	offset := 0

	_, err := parseName(message, &offset)

	assert.Error(t, err)
}
