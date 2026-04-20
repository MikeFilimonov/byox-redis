package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	terminator = "\r\n"
	arrayMark  = '*'
	stringMark = '$'
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	//
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {

		clientConnection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func() {
			replyToClient(clientConnection)
		}()
	}
}

func encodeBulkString(input string) string {
	return fmt.Sprintf("$%d%s%s%s", len(input), terminator, input, terminator)
}

func parseBulkString(input string) (string, string) {

	if len(input) == 0 || input[0] != stringMark {
		return "", input
	}

	// first line end
	endIdx := strings.Index(input, terminator)
	if endIdx == -1 {
		return "", input
	}

	rawLength := input[1:endIdx]
	length := 0
	_, err := fmt.Sscanf(rawLength, "%d", &length)
	if err != nil {
		return "", input
	}

	dataStartPosition := endIdx + 2
	dataEndPosition := dataStartPosition + length
	if dataEndPosition > len(input) {
		return "", input
	}

	remainder := input[dataEndPosition+len(terminator):]

	return input[dataStartPosition:dataEndPosition], remainder

}

func parseArray(input string) ([]string, string) {

	if len(input) == 0 || input[0] != arrayMark {
		return nil, input
	}

	// first line end
	firstLineEnd := strings.Index(input, terminator)
	if firstLineEnd == -1 {
		return nil, input
	}

	elementsCountRaw := input[1:firstLineEnd]

	elementsCount := 0
	_, err := fmt.Sscanf(elementsCountRaw, "%d", &elementsCount)
	if err != nil {
		return nil, input
	}

	remainder := input[firstLineEnd+2:]
	elements := make([]string, 0, elementsCount)

	for i := 0; i < elementsCount; i++ {
		if len(remainder) == 0 || remainder[0] != stringMark {
			break
		}

		bulkString, newRemainder := parseBulkString(remainder)
		elements = append(elements, bulkString)
		remainder = newRemainder

	}

	return elements, remainder

}

func replyToClient(clientConnection net.Conn) {

	for {

		buffer := make([]byte, 1024)
		n, err := clientConnection.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read data from client request: ", err.Error())
			os.Exit(1)
		}

		input := string(buffer[:n])
		reply := handleInput(input)
		clientConnection.Write([]byte(reply))
	}

}

func handleInput(input string) string {

	if len(input) == 0 {
		return ""
	}

	if input[0] == arrayMark {
		elements, _ := parseArray(input)
		if len(elements) == 0 {
			return ""
		}

		command := strings.ToLower(elements[0])
		switch command {
		case "ping":
			return "+PONG\r\n"
		case "echo":
			if len(elements) > 1 {
				return encodeBulkString(elements[1])
			}
			return "$0\r\n\r\n"
		default:
			return ""
		}

	}

	return ""
}
