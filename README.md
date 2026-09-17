 # Custom Redis-Compatible TCP Server in Go

A multi-threaded, custom Redis-compatible in-memory key-value store built from scratch in Go, interacting directly with raw TCP sockets.

## Features Implemented
- **Raw TCP Server:** Handles concurrent clients using Go routines and network sockets.
- **Thread-Safe Architecture:** Utilizes `sync.RWMutex` for safe concurrent map reads and writes.
- **Core Commands:**
  - `PING` / `ECHO`: Basic connection testing and echoing.
  - `SET` / `GET`: In-memory data storage and retrieval.
  - `EX` (TTL): Automatic time-to-live expiration and lazy cleanup.
  - `DEL`: Manual key deletion with RESP integer feedback.
  - `EXISTS`: Key presence verification.
  - `INCR`: Atomic numeric counters with type validation.

## How to Run
1. Start the server:
      bash
   go run main.go 

2.  Connect using `redis-cli` (or any Telnet-compatible client):
      bash
   redis-cli -p 6379
   
3.  Test the commands:
   ```redis
   PING
   ECHO hi there
   SET greeting "hello"
   GET greeting
   SET countdown 5 EX 10
   INCR countdown
   EXISTS countdown
   DEL greeting
   
