package main

type Message struct {
	Header     Header
	Question   Question
	Answer     ResourceRecord
	Authority  ResourceRecord
	Additional ResourceRecord
}
