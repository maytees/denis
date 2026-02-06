package main

import "strings"

type RecordType uint16

const (
	TypeA     RecordType = 1
	TypeNS    RecordType = 2
	TypeMD    RecordType = 3
	TypeMF    RecordType = 4
	TypeCNAME RecordType = 5
	TypeSOA   RecordType = 6
	TypeMB    RecordType = 7
	TypeMG    RecordType = 8
	TypeMR    RecordType = 9
	TypeNULL  RecordType = 10
	TypeWKS   RecordType = 11
	TypePTR   RecordType = 12
	TypeHINFO RecordType = 13
	TypeMINFO RecordType = 14
	TypeMX    RecordType = 15
	TypeTXT   RecordType = 16

	// Qtypes
	QTypeAXFR  RecordType = 252
	QTypeMAILB RecordType = 253
	QTypeMAILA RecordType = 254
	QTypeAll   RecordType = 255
)

var RecordTypeNames = map[RecordType]string{
	TypeA:     "A",
	TypeNS:    "NS",
	TypeMD:    "MD",
	TypeMF:    "MF",
	TypeCNAME: "CNAME",
	TypeSOA:   "SOA",
	TypeMB:    "MB",
	TypeMG:    "MG",
	TypeMR:    "MR",
	TypeNULL:  "NULL",
	TypeWKS:   "WKS",
	TypePTR:   "PTR",
	TypeHINFO: "HINFO",
	TypeMINFO: "MINFO",
	TypeMX:    "MX",
	TypeTXT:   "TXT",

	// QTypes
	QTypeAXFR:  "AXFR",
	QTypeMAILB: "MAILB",
	QTypeMAILA: "MAILA",
	QTypeAll:   "*",
}

type ClassValue uint16

const (
	ClassIN ClassValue = 1
	ClassCS ClassValue = 2
	ClassCH ClassValue = 3
	ClassHS ClassValue = 4

	// QClass values
	QClassAny ClassValue = 255
)

var ClassValueNames = map[ClassValue]string{
	ClassIN: "IN",
	ClassCS: "CS",
	ClassCH: "CH",
	ClassHS: "HS",

	// QClass values
	QClassAny: "*",
}

type ResourceRecord struct {
	Name     Name
	Type     RecordType
	Class    ClassValue
	TTL      uint32
	RDlength uint16

	// Should this be composite with raw?
	RData string
}

// TODO: Slow linear search, optimize later
func findRecordByName(records []Record, target string) (Record, bool) {
	for _, record := range records {
		// Records are always stored lowercase
		if record.Name == strings.ToLower(target) {
			return record, true
		}
	}

	return Record{}, false
}
