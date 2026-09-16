package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// bind to port 6379 (standard redis port)

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("failed to bind to port 6379: %v\n", err)
		return
	}

	defer listener.Close()

	fmt.Println("custom redis server listening on port 6379...")

	// continous accept loop for incoming clients

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		// handle each client concurrently using a goroutine
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	// read/write loop for the individual client connection

	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		// convert incoming bytes to a clean string
		request := string(buf[:n])

		// check what command the client sent (ignoring case and whitespace)

		upperReq := strings.ToUpper(strings.TrimSpace(request))

		if strings.Contains(upperReq, "PING") {
			_, _ = conn.Write([]byte("+PONG\r\n"))
		} else if strings.HasPrefix(upperReq, "ECHO") {

			// if its an ECHO command, let's pull out the message part
			// (for now , we'll echo back whatever followed ECHO)

			parts := strings.SplitN(request, " ", 2)
			if len(parts) > 1 {
				message := strings.TrimSpace(parts[1])

				// respond with a RESP bulk string or simple string

				_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of argumentsfor 'ECHO' command\r\n"))
			}
		} else {
			_, _ = conn.Write([]byte("+OK\r\n"))

		}
	}

}
