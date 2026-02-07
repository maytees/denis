package dns

import "fmt"

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
