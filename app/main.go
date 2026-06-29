package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main() {	
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
			break
		}	
	resp := strings.Replace(string(buff[:n]), "\r\n", "", -1)
		if strings.Contains(resp, "PING") {
			conn.Write([]byte("+PONG\r\n"))
		}
	
}
}