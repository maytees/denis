package main

import "fmt"

// For testing
const DEV_PORT = 5353
const PROD_PORT = 53

func cliWelcome() {
	fmt.Println("Welcome to the DENIS DNS server!")
}

type ServerConfig struct {
	Port int
}

func (c ServerConfig) PrintConfig() {
	fmt.Printf("Port: %v\n", c.Port)
}

func main() {
	cliWelcome()

	config := ServerConfig{
		Port: 53,
	}

	config.PrintConfig()
}
