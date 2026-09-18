package main

import (
	"bufio"
	"fmt"
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
	fmt.Println("redis server listenin on 127.0.0.1:6379...")
	for {

		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection:", err)
			continue
		}
		go s.handleClient(conn)

	}

}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		line := scanner.Text()
		parts := strings.Fields(line)
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

		default:
			_, _ = conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))

		}

	}

}
