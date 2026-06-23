package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	terminator        = "\r\n"
	newLiner          = '\n'
	arrayMark         = '*'
	bulkStringMark    = '$'
	stringMark        = '+'
	timeTolive        = 3600 * time.Second // seconds in an hour
	defaultParseError = "failed to read the input: "
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var storage = &Storage{data: make(map[string]*Entry)}

func main() {

	// get the item count
	// run the loop until done for collecting input parts

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

func parseBulkString(reader *bufio.Reader) (string, error) {

	input, err := reader.ReadString(newLiner)
	if err != nil {
		fmt.Println(defaultParseError, err.Error())
		return "", err
	}

	input = strings.TrimSuffix(input, terminator)

	if len(input) == 0 || input[0] != bulkStringMark {
		return "", fmt.Errorf("expected bulk string mark, got %q", input)
	}

	length, err := strconv.Atoi(input[1:])
	if err != nil {
		fmt.Println("failed to parse bulk string length: ", err.Error())
		return "", err
	}

	buffer := make([]byte, length+len(terminator))

	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		fmt.Println("failed to read array elements")
		return "", err
	}

	return string(buffer[:length]), nil

}

func parseArray(reader *bufio.Reader) ([]string, error) {

	input, err := reader.ReadString(newLiner)
	if err != nil {
		fmt.Println(defaultParseError, err.Error())
		return []string{}, err
	}

	input = strings.TrimSuffix(input, terminator)
	length, err := strconv.Atoi(input[1:])
	if err != nil {
		fmt.Println("failed to parse array length: ", err.Error())
		return []string{}, err
	}

	elements := make([]string, 0, length)

	for i := 0; i < length; i++ {

		element, err := parseBulkString(reader)
		if err != nil {
			fmt.Println(defaultParseError, err.Error())
			return []string{}, err
		}
		elements = append(elements, element)

	}

	return elements, nil

}

func replyToClient(clientConnection net.Conn) {

	defer clientConnection.Close()
	reader := bufio.NewReader(clientConnection)

	for {

		reply, err := handleInput(reader)
		if err != nil {
			fmt.Println("Failed to read data from client request: ", err.Error())
			return
		}
		if _, err := clientConnection.Write([]byte(reply)); err != nil {
			fmt.Println("Failed to send the reply: ", err.Error())
			return
		}
	}

}

func handleInput(reader *bufio.Reader) (string, error) {

	firstItem, err := reader.ReadByte()
	if err != nil {
		fmt.Println(defaultParseError, err.Error())
		return "", err
	}

	defaultReply := fmt.Sprintf("0%s%s", terminator, terminator)

	if rune(firstItem) != arrayMark {
		return defaultReply, nil
	}

	reader.UnreadByte() // parseArray gotta catch'em all

	elements, err := parseArray(reader)
	if err != nil {
		return "", err
	}
	if len(elements) == 0 {
		return defaultReply, nil
	}

	command := strings.ToLower(elements[0])
	switch command {
	case "ping":
		return fmt.Sprintf("%vPONG%s", string(stringMark), terminator), nil
	case "echo":
		if len(elements) > 1 {
			return encodeBulkString(elements[1]), nil
		}
		return defaultReply, nil
	case "get":
		return getValue(elements, defaultReply), nil
	case "set":
		return setValue(elements, defaultReply), nil
	default:
		return defaultReply, nil
	}

}

func getValue(elements []string, defaultReply string) string {

	if len(elements) > 1 {

		data, ok := storage.Get(elements[1])
		if !ok {
			return fmt.Sprintf("%v-1%s", string(bulkStringMark), terminator)
		}

		return encodeBulkString(fmt.Sprintf("%v", data.Value))

	}
	return defaultReply

}

func setValue(elements []string, defaultReply string) string {

	if len(elements) >= 3 {

		key, value := elements[1], elements[2]

		data := &Entry{
			Value: value,
		}

		storage.Set(key, data, timeTolive)

		if len(elements) == 5 {

			opt := strings.ToLower(elements[3])
			if opt != "px" && opt != "ex" {
				return defaultReply
			}

			val, err := strconv.Atoi(elements[4])
			if err != nil {
				return defaultReply
			}

			var lifespan time.Duration
			if opt == "px" {
				lifespan = time.Duration(val) * time.Millisecond
			} else {
				lifespan = time.Duration(val) * time.Second
			}
			storage.Set(key, data, lifespan)

		}

		return fmt.Sprintf("%sOK%s", string(stringMark), terminator)
	}
	return defaultReply

}
