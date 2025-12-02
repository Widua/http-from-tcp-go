package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/widua/http-from-tcp-go/internal/request"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)
	go func() {
		defer close(strChan)
		defer f.Close()

		message := ""

		buff := make([]byte, 8)

		for {
			_, err := f.Read(buff)
			if err == io.EOF {
				break
			}
			strBuf := string(buff)
			message += strBuf
			if strings.Contains(strBuf, "\n") {
				lines := strings.Split(message, "\n")
				strChan <- lines[0]
				message = strings.Join(lines[1:], "\n")
				continue
			}

		}
	}()
	return strChan
}

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		defer conn.Close()
		fmt.Println("Connection accepted...")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error: %v", err)
			continue
		}
		fmt.Printf("Request line: \n - Method: %v\n- Target: %v\n- Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println(req.Headers)
		fmt.Printf("Body:\n%v", string(req.Body))

		fmt.Printf("Connection closed...")
	}

}
