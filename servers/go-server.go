/*
To test the slow loris I thought it would be useful to build my own
http server.
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	httpRequest := ""
	for {
		n, err := conn.Read(buffer)
		paritalRequest := string(buffer[:n])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		// If this needed to be fast you should use a slice of bytes so that
		// you don't need to creat a new string each time you append to httpRequest.
		httpRequest += paritalRequest
		fmt.Printf("Request so far: %q \n", httpRequest)
		fmt.Printf("End values: %q\n", httpRequest[len(httpRequest)-4:])
		if httpRequest[len(httpRequest)-4:] == "\r\n\r\n" {
			break
		}
	}

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"Connection: close\r\n" +
		"Server-Response: close\r\n\r\n" +
		// I got confused here we ended the headers with \r\n\r\n
		// and proceeded to write the response body, but obviously
		// GET responses needed to contain a body.
		"server response!"

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

func server(port string) (err error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	defer func() { err = errors.Join(listener.Close(), err) }()

	fmt.Println("Server is listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// might be worth thinking about better ways of handling errors
		// instead of relying on goroutines printing to stdout.
		go handleConnection(conn)
	}
}

func main() {
	port := "8080"
	err := server(port)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
