package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	listener net.Listener
	store    *Store
}

func NewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
		store:    NewStore(),
	}, nil
}

func (s *Server) Start() error {
	fmt.Println("Custom Redis server listening on 127.0.0.1:6379...")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handleClient(conn)
	}
}

// parseCommand reads either a strict RESP array or falls back to plain text.
func parseCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	if len(line) == 0 {
		return nil, nil
	}

	// If it starts with '*', it's a strict RESP array from redis-cli
	if line[0] == '*' {
		var numElements int
		_, err := fmt.Sscanf(line, "*%d", &numElements)
		if err != nil {
			return nil, err
		}

		var args []string
		for i := 0; i < numElements; i++ {
			// Read bulk string length line (e.g. $3)
			lenLine, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			lenLine = strings.TrimRight(lenLine, "\r\n")
			if len(lenLine) == 0 || lenLine[0] != '$' {
				return nil, fmt.Errorf("expected bulk string header")
			}

			var strLen int
			_, err = fmt.Sscanf(lenLine, "$%d", &strLen)
			if err != nil {
				return nil, err
			}

			// Read the actual string payload + \r\n
			buf := make([]byte, strLen+2)
			_, err = io.ReadFull(reader, buf)
			if err != nil {
				return nil, err
			}
			args = append(args, string(buf[:strLen]))
		}
		return args, nil
	}

	// Fallback to plain text space-separated commands (for manual testing/scripts)
	return strings.Fields(line), nil
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		parts, err := parseCommand(reader)
		if err != nil {
			if err != io.EOF {
				// Handle read errors if needed
			}
			break
		}
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "PING":
			_, _ = conn.Write([]byte("+PONG\r\n"))

		case "ECHO":
			if len(parts) >= 2 {
				msg := strings.Join(parts[1:], " ")
				_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'echo'\r\n"))
			}

		case "SET":
			if len(parts) >= 3 {
				key := parts[1]
				val := parts[2]
				var expiresAt time.Time

				// Check for EX modifier
				if len(parts) >= 5 && strings.ToUpper(parts[3]) == "EX" {
					var seconds int
					_, err := fmt.Sscanf(parts[4], "%d", &seconds)
					if err == nil {
						expiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
					}
				}

				s.store.Set(key, val, expiresAt)
				_, _ = conn.Write([]byte("+OK\r\n"))
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'set'\r\n"))
			}

		case "GET":
			if len(parts) >= 2 {
				key := parts[1]
				val, found := s.store.Get(key)
				if !found {
					_, _ = conn.Write([]byte("$-1\r\n"))
				} else {
					_, _ = conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'get'\r\n"))
			}

		case "DEL":
			if len(parts) >= 2 {
				key := parts[1]
				deleted := s.store.Delete(key)
				if deleted {
					_, _ = conn.Write([]byte(":1\r\n"))
				} else {
					_, _ = conn.Write([]byte(":0\r\n"))
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'del'\r\n"))
			}

		case "EXISTS":
			if len(parts) >= 2 {
				key := parts[1]
				if s.store.Exists(key) {
					_, _ = conn.Write([]byte(":1\r\n"))
				} else {
					_, _ = conn.Write([]byte(":0\r\n"))
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'exists'\r\n"))
			}

		case "EXPIRE":
			if len(parts) >= 3 {
				key := parts[1]
				var seconds int
				_, err := fmt.Sscanf(parts[2], "%d", &seconds)
				if err != nil {
					_, _ = conn.Write([]byte("-ERR value is not an integer or out of range\r\n"))
					break
				}

				success := s.store.Expire(key, time.Duration(seconds)*time.Second)
				if success {
					_, _ = conn.Write([]byte(":1\r\n"))
				} else {
					_, _ = conn.Write([]byte(":0\r\n"))
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'expire'\r\n"))
			}

		case "INCR":
			if len(parts) >= 2 {
				key := parts[1]
				newVal, err := s.store.Incr(key)
				if err != nil {
					_, _ = conn.Write([]byte("-ERR value is not an integer or out of range\r\n"))
				} else {
					_, _ = conn.Write([]byte(fmt.Sprintf(":%d\r\n", newVal)))
				}
			} else {
				_, _ = conn.Write([]byte("-ERR wrong number of arguments for 'incr'\r\n"))
			}

		default:
			_, _ = conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))
		}
	}
}
