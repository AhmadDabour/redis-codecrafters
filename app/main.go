package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
	"strconv"
)

var _ = net.Listen
var _ = os.Exit

func main() {	
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn)  {
	buff := make([]byte, 1024)
	for {
		n, err := c.Read(buff)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
			break
		}
		if strings.Contains(strings.ToLower(string(buff)), "echo") { 
			go echoParser(string(buff[:n]), c)
		}
		resp := strings.Replace(string(buff[:n]), "\r\n", "", -1)
		if strings.Contains(resp, "PING") {
			c.Write([]byte("+PONG\r\n"))
			continue
		}
	}
}

func echoParser(buff string, c net.Conn) {
	reader := bufio.NewReader(strings.NewReader(buff))
	for i := 0; i < 14; i++ {
		reader.ReadByte()
	}
	
	// s, _ := reader.Peek(1)
	// if  s[0] == '*' {
	// 	reader.ReadByte()
	// 	reader.ReadByte()
	// }
	b, _ := reader.ReadByte()
	if b != '$' {
		fmt.Println("Invalid type")
		os.Exit(1)
	}
	size, _ := reader.ReadByte()
	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	reader.ReadByte()
	reader.ReadByte()

	resp := make([]byte, strSize)
	reader.Read(resp)
	c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n",strSize, string(resp))))
}