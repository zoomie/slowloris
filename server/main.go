package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the request
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// Convert the request to a string and print it
	request := string(buffer[:n])
	fmt.Println("Received request:")
	fmt.Println(request)

	// Basic check to ensure it's a GET request
	if strings.HasPrefix(request, "GET") {
		// Create a basic HTTP response
		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: 13\r\n" +
			"Connection: close\r\n\r\n" +
			"Hello, world!"

		// Write the response to the connection
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	} else {
		// Respond with 400 Bad Request for non-GET requests
		response := "HTTP/1.1 400 Bad Request\r\n" +
			"Connection: close\r\n\r\n"
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}

func main() {
	// Define the port to listen on
	port := "8080"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", port)

	for {
		// Accept a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}
