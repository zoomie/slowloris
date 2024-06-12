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
		"X-Header: new\r\n"

	if _, err = conn.Write([]byte(startRequest)); err != nil {
		return fmt.Errorf("error writing to connection: %w", err)
	}

	for range 5 {
		header := "UserAgent: this is random header\r\n"
		if _, err := conn.Write([]byte(header)); err != nil {
			return fmt.Errorf("error writing to connection: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	endHttp := "\r\n"
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
	rsp := ""
	var err error
	for {
		numberBytesRead := 0
		numberBytesRead, err = conn.Read(buffer)
		fmt.Println(numberBytesRead)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading from connection: %w", err)
		}
		rsp += string(buffer[:numberBytesRead])
	}
	return rsp, nil
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
