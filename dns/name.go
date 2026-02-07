package dns

import (
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

func parseName(message []byte, offset *int) Name {
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

		builder.WriteString(string(message[*offset:(*offset + length)]))
		*offset += length
	}

	return Name{
		Raw:    message[start:*offset],
		Domain: builder.String(),
	}
}
