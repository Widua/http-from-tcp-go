package main

import (
	"fmt"
	"net"
)

func main() {
	endpoint, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	conn, err := net.DialUDP("udp")
}
