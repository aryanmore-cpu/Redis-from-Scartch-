package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

//  in-memory store and a mutex to make it thread-safe for concurrent clients

var (
	kvStore = make(map[string]string)
	mu      sync.RWMutex
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

		request := string(buf[:n])

		parts := strings.Fields(request) // splits by whitespace
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {

		case "PING":
			_, _ = conn.Write([]byte("+PONG\r\n"))

		case "ECHO":

			if len(parts) > 1 {
				msg := parts[1]
				_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'echo'\r\n"))
			}

		case "SET":
			if len(parts) > 3 {
				key := parts[1]
				val := parts[2]

				mu.Lock()
				kvStore[key] = val
				mu.Unlock()

				_, _ = conn.Write([]byte("+OK\r\n"))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'SET'\r\n"))
			}

		case "GET":
			if len(parts) >= 2 {
				key := parts[1]

				mu.RLock()
				val, exists := kvStore[key]
				mu.RUnlock()

				if exists {
					_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
				} else {
					_, _ = conn.Write([]byte("-$-1\r\n")) // redis null bulk string or missing keys
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'GET'\r\n"))
			}

		default:
			_, _ = conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}

}
