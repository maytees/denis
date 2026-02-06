package main

import "fmt"

type Message struct {
	Header     Header
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
