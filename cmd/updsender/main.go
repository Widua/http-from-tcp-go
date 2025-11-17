package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	endpoint, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, endpoint)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		conn.Write([]byte(str))
	}
	defer conn.Close()
}
