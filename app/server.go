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
	requestLine := strings.SplitN(string(buf), "\n", 2)[0]
	if strings.HasPrefix(requestLine, "GET /user-agent HTTP/1.1") {
		// Handle user-agent request
		res := strings.Split(string(buf), "\n")
		userAgent := "User-Agent not found"
		for _, line := range res {
			if strings.HasPrefix(line, "User-Agent:") {
				parts := strings.SplitN(line, ": ", 2)
				if len(parts) == 2 {
					userAgent = strings.TrimSpace(parts[1])
				}
				break
			}
		}
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		conn.Write([]byte(resp))
	} else if strings.HasPrefix(requestLine, "GET /echo/") {
		// Handle echo request
		message := strings.TrimPrefix(requestLine, "GET /echo/")
		message = strings.TrimSuffix(message, " HTTP/1.1")
		message = strings.TrimSpace(message)
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
		conn.Write([]byte(resp))
	} else if strings.HasPrefix(requestLine, "GET / HTTP/1.1") {
		// Respond to the root request
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		// Respond with 404 for other requests
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
