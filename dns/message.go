package dns

import (
	"fmt"
)

type Message struct {
	Header Header
	// Don't do multiple questions for now.
	Question   Question
	Answer     []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
	Raw        []byte
}

func (m Message) String() string {
	return fmt.Sprintf("Header: %v\nQuestion: %v\nAnswer(s): %v\nAuthority(s): %v\nAdditional(s): %v",
		m.Header,
		m.Question,
		m.Answer,
		m.Authority,
		m.Additional,
	)
}

func parseMessage(message []byte) (*Message, int, error) {
	offset := 12
	header := parseHeader(message[:offset])

	questions, err := parseQuestions(message, &offset, header.QDCount)
	if err != nil {
		return nil, 0, err
	}

	return &Message{
		Header: header,
		// Only do single question for now.
		// See struct in question.go
		Question: questions[0],
		Raw:      message,
	}, offset, nil

	// NOTE: do smthn i forgot what
	// TODO: handle authority & additional if necessary

	// TODO: Get by record & name
	// record, ok := config.FindRecordByName(*records, name.Domain)
	// if !ok {
	// 	forwardQuery(dnsConfig.Upstream, message, connection, clientAddr)
	// 	return
	// }

	// _, port, variant, err := util.ParseAddress(record.Value)

	// if err != nil {
	// 	fmt.Println("Could not parse address: ", err)
	// 	return
	// }

	// if variant != "v4" {
	// 	fmt.Println("Record should be an IPv4 address!")
	// 	return
	// }

	// if port != "" {
	// 	fmt.Println("Record should not have port number!")
	// 	return
	// }

	// return &Message{
	// 	Header:     header,
	// 	Question:   nil,
	// 	Answer:     nil,
	// 	Authority:  nil,
	// 	Additional: nil,
	// 	Raw:        message,
	// }

	// return nil, nil
}
