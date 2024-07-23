package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func init() {
	flag.StringVar(&directory, "directory", "", "Specify the directory to use")
}

// The main function is the entry point of the program.
// It sets up a TCP listener on port 4221 and handles incoming connections.
func main() {
	flag.Parse()
	fmt.Println("Logs from your program will appear here!")

	// Listen for incoming TCP connections on all available network interfaces using port 4221.
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		// If binding to the port fails, print an error message and exit the program.
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Server is listening on port 4221")

	// Accept a new incoming connection.
	for {
		conn, err := l.Accept()
		if err != nil {
			// If accepting the connection fails,
			//print an error message and continue.
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		// Handle the connection concurrently using a Goroutine.
		go handleConnection(conn)
	}

}

// The handleConnection function handles an individual connection.
func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function completes.

	// Create a buffer to store the incoming data.
	buf := make([]byte, 1024)
	// Read the incoming data into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		// If reading fails, exit the function.
		return
	}

	// Parse the first line of the request to determine the handling logic.
	requestLine := strings.SplitN(string(buf), "\n", 2)[0]
	// Handle the request based on the request line.
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
		// Construct and send the response containing the user-agent information.
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		conn.Write([]byte(resp))
	} else if strings.HasPrefix(requestLine, "GET /echo/") {
		// Handle echo request
		message := strings.TrimPrefix(requestLine, "GET /echo/")
		message = strings.TrimSpace(message)
		message = strings.TrimSuffix(message, " HTTP/1.1")
		// Construct and send the response containing the echoed message.
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
		conn.Write([]byte(resp))
	} else if strings.HasPrefix(requestLine, "GET / HTTP/1.1") {
		// Respond to the root request
		// Send a simple HTTP 200 OK response for the root path.
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(requestLine, "GET /files/") {
		fileName := strings.TrimPrefix(requestLine, "GET /files/")
		fileName = strings.TrimSpace(fileName)
		fileName = strings.TrimSuffix(fileName, " HTTP/1.1")

		filePath := directory + fileName
		fmt.Println("file:", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)
		conn.Write([]byte(resp))
	} else {
		// Respond with 404 for other requests
		// Send a HTTP 404 Not Found response for unhandled paths.
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
