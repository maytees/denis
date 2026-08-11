// NOTE: This test file was written by AI (Claude).
package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Every defined RecordType constant should have a display name
func TestRecordTypeNamesComplete(t *testing.T) {
	types := []RecordType{
		TypeA, TypeNS, TypeMD, TypeMF, TypeCNAME, TypeSOA, TypeMB, TypeMG,
		TypeMR, TypeNULL, TypeWKS, TypePTR, TypeHINFO, TypeMINFO, TypeMX,
		TypeTXT, QTypeAXFR, QTypeMAILB, QTypeMAILA, QTypeAll,
	}

	for _, rt := range types {
		name, ok := RecordTypeNames[rt]
		assert.True(t, ok, "RecordType %d has no name", rt)
		assert.NotEmpty(t, name)
	}
}

func TestClassValueNamesComplete(t *testing.T) {
	classes := []ClassValue{ClassIN, ClassCS, ClassCH, ClassHS, QClassAny}

	for _, cv := range classes {
		name, ok := ClassValueNames[cv]
		assert.True(t, ok, "ClassValue %d has no name", cv)
		assert.NotEmpty(t, name)
	}
}

func TestResourceRecordString(t *testing.T) {
	rr := ResourceRecord{
		Name:     Name{Domain: "my.google"},
		Type:     TypeA,
		Class:    ClassIN,
		TTL:      300,
		RDlength: 4,
		RData:    "142.251.16.138",
	}

	s := rr.String()

	assert.Contains(t, s, "my.google")
	assert.Contains(t, s, "A")
	assert.Contains(t, s, "IN")
	assert.Contains(t, s, "142.251.16.138")
}

func TestQuestionString(t *testing.T) {
	q := Question{
		QName:  Name{Domain: "my.google"},
		QType:  TypeA,
		QClass: ClassIN,
	}

	s := q.String()

	assert.Contains(t, s, "my.google")
	assert.Contains(t, s, "A")
	assert.Contains(t, s, "IN")
}
