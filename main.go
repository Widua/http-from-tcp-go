package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	messages, _ := os.Open("messages.txt")
	message := ""

	for {
		buff := make([]byte, 8)
		_, err := messages.Read(buff)

		if err == io.EOF {
			break
		}
		strBuf := string(buff)
		message += strBuf
		if strings.Contains(strBuf, "\n") {
			lines := strings.Split(message, "\n")
			fmt.Printf("read: %s\n", lines[0])
			message = strings.Join(lines[1:], "\n")
			continue
		}
	}

}
