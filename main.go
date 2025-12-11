package main

import (
	"fmt"
	"log"
	"net"
)

func cliWelcome() {
	fmt.Println("Welcome to the DENIS DNS server!")
}

type ServerConfig struct {
	Port int
}

func main() {
	cliWelcome()

	udpAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:5354")
	if err != nil {
		log.Fatalln(err)
	}

	connection, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()

	// UDP messages 512 octets (bytes) or less
	buffer := make([]byte, 512)

	for {
		input, clientAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from packet %v\n", err)
			continue
		}

		data := buffer[:input]
		log.Printf("FROM \"%v\" (%d bytes)\n%x", clientAddr.String(), input, data)
	}
}
