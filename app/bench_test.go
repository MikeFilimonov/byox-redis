package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	protocol = "tcp"
	server   = "127.0.0.1:0"
	getCount = 1000
	setCount = 1000
)

func newDummyServer(t *testing.T) string {

	t.Helper()
	listener, err := net.Listen(protocol, server)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			go replyToClient(connection)
		}
	}()
	return listener.Addr().String()

}

func TestBenchmark(t *testing.T) {

	address := newDummyServer(t)

	connection, err := net.Dial(protocol, address)
	if err != nil {
		t.Fatal(err)
	}

	defer connection.Close()

	reader := bufio.NewReader(connection)

	send := func(args ...string) {

		var sb strings.Builder
		fmt.Fprintf(&sb, "%c%d%s", arrayMark, len(args), terminator)
		for _, a := range args {
			fmt.Fprintf(&sb, "%c%d%s%s%s",
				bulkStringMark, len(a), terminator, a, terminator)
		}
		connection.Write([]byte(sb.String()))
	}

	readReply := func() {
		line, _ := reader.ReadString(newLiner)
		if len(line) > 0 && line[0] == bulkStringMark {
			reader.ReadString(newLiner)
		}
	}

	startSetting := time.Now()
	for i := 0; i < setCount; i++ {
		send("SET", fmt.Sprintf("key:%d", i), fmt.Sprintf("val:%d", i))
		readReply()
	}

	setDuration := time.Since(startSetting)

	startGetting := time.Now()
	for i := 0; i < getCount; i++ {
		send("GET", fmt.Sprintf("key:%d", i))
		readReply()
	}

	getDuration := time.Since(startGetting)

	t.Logf("SET %d ops in %v (%.0f ops/sec)", setCount, setDuration,
		float64(setCount)/setDuration.Seconds())
	t.Logf("GET %d ops in %v (%.0f ops/sec)", getCount, getDuration,
		float64(getCount)/getDuration.Seconds())

}
