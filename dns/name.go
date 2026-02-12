package dns

import (
	"errors"
	"fmt"
	"strings"
)

type Name struct {
	Raw    []byte // Includes label lengths
	Domain string
}

func (n Name) String() string {
	return fmt.Sprintf("%s", n.Domain)
}

func parseName(message []byte, offset *int) (*Name, error) {
	start := *offset
	var builder strings.Builder

	for {
		length := int(message[*offset])
		*offset += 1

		if length == 0 {
			break
		}

		if builder.Len() != 0 {
			builder.WriteByte('.')
		}

		jump := *offset + length
		if jump > len(message) {
			return nil, errors.New("Message too short")
		}

		_, err := builder.WriteString(strings.ToLower(string(message[*offset:jump])))
		if err != nil {
			return nil, err
		}
		*offset += length
	}

	return &Name{
		Raw:    message[start:*offset],
		Domain: builder.String(),
	}, nil
}
