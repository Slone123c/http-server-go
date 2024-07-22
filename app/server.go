package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return
	}
	if strings.HasPrefix(string(buf), "GET /user-agent HTTP/1.1") {
		res := strings.Split(string(buf), "\n")
		usetAgent := strings.Split(res[3], ":")[1]
		resp := fmt.Sprintf("HTTP/1."+
			"1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n"+
			"\r\n%s", len(usetAgent), usetAgent)
		conn.Write([]byte(resp))
	} else if strings.HasPrefix(string(buf), "GET / HTTP/1.1") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
