package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

func oneSlowLorisCall(server string, port string) (err error) {
	addr, err := net.ResolveTCPAddr("tcp", server+":"+port)
	if err != nil {
		return fmt.Errorf("error resolving address: %v", err)
	}

	// Need to create TCP seperately to control so we can control
	// the speed of data written to the socket.
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return fmt.Errorf("error setting up tcp connection: %w", err)
	}
	defer func() { err = errors.Join(conn.Close(), err) }()

	startRequest := "GET / HTTP/1.1\r\n" +
		"Host: " + server + "\r\n" +
		"Connection: close\r\n"

	if _, err = conn.Write([]byte(startRequest)); err != nil {
		return fmt.Errorf("error writing to connection: %w", err)
	}

	userAgent := "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\r\n"

	for _, b := range []byte(userAgent) {
		if _, err := conn.Write([]byte{b}); err != nil {
			return fmt.Errorf("error writing to connection: %w", err)
		}
		time.Sleep(1 * time.Millisecond)
	}

	endHttp := "\r\n\r\n"
	if _, err := conn.Write([]byte(endHttp)); err != nil {
		return fmt.Errorf("error writing to connection: %w", err)
	}

	rsp, err := readResponse(conn)
	if err != nil {
		return fmt.Errorf("error reading from connection: %w", err)
	}
	fmt.Println(rsp)

	return nil
}

// readResponse can be used to test the server response
func readResponse(conn net.Conn) (string, error) {
	buffer := make([]byte, 4096)
	numberBytesRead := 0
	var err error
	for {
		numberBytesRead, err = conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading from connection: %w", err)
		}
	}
	return string(buffer[:numberBytesRead]), nil
}

func main() {
	server := "localhost"
	port := "8000"
	err := oneSlowLorisCall(server, port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
