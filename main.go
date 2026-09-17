package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// item holds the value and an optional expiration time

type Item struct {
	value     string
	expiresAt time.Time // zero time means no expiration
}

//  in-memory store and a mutex to make it thread-safe for concurrent clients

var (
	kvStore = make(map[string]Item)
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
				var expiresAt time.Time

				// check if an expiration argument was passed (eg SET key val EX 20 )

				if len(parts) >= 5 && strings.ToUpper(parts[3]) == "EX" {
					seconds, err := strconv.Atoi(parts[4])
					if err == nil {
						expiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
					}

				}

				mu.Lock()
				kvStore[key] = Item{
					value:     val,
					expiresAt: expiresAt,
				}
				mu.Unlock()

				_, _ = conn.Write([]byte("+OK\r\n"))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'SET'\r\n"))
			}

		case "GET":
			if len(parts) >= 2 {
				key := parts[1]

				mu.Lock()
				item, exists := kvStore[key]

				// check if the key exists AND if it has expired

				if !exists {
					mu.Unlock()
					_, _ = conn.Write([]byte("-$-1\r\n"))
					break
				}

				if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {

					// its expired!! delete it from the map and treat it as missing

					delete(kvStore, key)
					mu.Unlock()
					_, _ = conn.Write([]byte("$-1\r\n"))
					break
				}

				val := item.value
				mu.Unlock()
				_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'GET'\r\n"))
			}
		default:
			_, _ = conn.Write([]byte("-ERR unknown command\r\n"))

		}
	}
}
