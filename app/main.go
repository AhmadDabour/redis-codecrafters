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
var data = make(map[string]string)

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
			fmt.Println("Error reading input: ")
			break
		}
		result := respParser(string(buff[:n]))
		switch result[0] {
			case "PING":
				c.Write([]byte("+PONG\r\n"))
			case "ECHO":
				c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(result[1]), result[1])))
			case "SET":
				data[result[1]] = result[2]
				c.Write([]byte("+OK\r\n"))
			case "GET":
				if res, ok := data[result[1]]; !ok {
					c.Write([]byte("$-1\r\n"))
				} else { 
				c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(res), res)))
			}
		}
	}
}

func respParser(buff string) []string{
	reader := bufio.NewReader(strings.NewReader(buff))
	s, _ := reader.ReadByte()
	if s != '*' {
		fmt.Println("Error parsing command")
		return nil
	}
	b, _ := reader.ReadByte() 
	sliceSize, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		fmt.Println("Error parsing int: ", err.Error())
	}
	input := []string{}

	for i := 0; i < int(sliceSize); i++ { 
	reader.ReadByte()
	reader.ReadByte()
	p, _ := reader.ReadByte()
	if p != '$' {
		fmt.Println("Error parsing command")
		return nil
	}
	size, _ := reader.ReadByte()
	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		fmt.Println("Error parsing int: ", err.Error())
		return nil
	}
	reader.ReadByte()
	reader.ReadByte()
	resp := make([]byte, strSize)
	reader.Read(resp)
	input = append(input, string(resp))
	}
	strings.ToLower(input[0])
	return input
}