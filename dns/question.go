package dns

import (
	"encoding/binary"
	"errors"
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
		name, err := parseName(message, offset)
		if err != nil {
			return nil, err
		}
		fmt.Println("Resolved domain:", name.Domain)

		qTypeJump := *offset + 2
		if qTypeJump > len(message) {
			return nil, errors.New("Message cuts off before QType")
		}
		qType := binary.BigEndian.Uint16(message[*offset:qTypeJump])
		*offset = qTypeJump

		qClassJump := *offset + 2
		if qClassJump > len(message) {
			return nil, errors.New("Message cuts off before QClass")
		}
		qClass := binary.BigEndian.Uint16(message[*offset:qClassJump])
		*offset = qClassJump

		questions = append(questions, Question{
			QName:  *name,
			QType:  RecordType(qType),
			QClass: ClassValue(qClass),
		})
	}

	return questions, nil
}
