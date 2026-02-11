package dns

import (
	"encoding/binary"
	"fmt"
)

type Question struct {
	QName  Name
	QType  RecordType
	QClass ClassValue
}

func (q Question) String() string {
	return fmt.Sprintf("Name: %v | QType: %s | QClass: %s",
		q.QName,
		RecordTypeNames[q.QType],
		ClassValueNames[q.QClass],
	)
}

func parseQuestions(message []byte, offset *int, qdCount uint16) ([]Question, error) {
	// Slice with size 0, but capacity qdCount
	questions := make([]Question, 0, qdCount)

	for range qdCount {
		name := parseName(message, offset)
		fmt.Println("Resolved domain:", name.Domain)

		qType := binary.BigEndian.Uint16(message[*offset : *offset+2])
		*offset += 2

		qClass := binary.BigEndian.Uint16(message[*offset : *offset+2])
		*offset += 2

		questions = append(questions, Question{
			QName:  name,
			QType:  RecordType(qType),
			QClass: ClassValue(qClass),
		})
	}

	return questions, nil
}
