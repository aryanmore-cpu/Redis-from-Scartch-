# Go-Redis: A Custom In-Memory Key-Value Store

A lightweight, multi-threaded, Redis-compatible TCP server built from absolute scratch in Go. This project implements custom network protocol handling, concurrent socket management, thread-safe memory stores, and atomic data structures without external framework dependencies.

---

## Features Implemented

* **Raw TCP Server & Concurrency:** Manages concurrent client connections using Go goroutines and underlying network sockets.
* **Thread-Safe Architecture:** Utilizes `sync.RWMutex` to ensure safe concurrent read/write access to the in-memory map store.
* **Custom Data Structures & Expiration Engine:** 
  * Implements custom `Item` structs containing values and `time.Time` timestamp tracking.
  * Features automatic time-to-live (`EX`/TTL) expiration with lazy evaluation and runtime cleanup checks.
* **Core Command Set:**
  * `PING` / `ECHO`: Connection verification and loopback echoing.
  * `SET` / `GET`: High-performance key-value storage and retrieval with precise argument parsing.
  * `EX` (TTL): Automatic time-to-live key expiration.
  * `DEL`: Manual key deletion returning RESP integer status responses (`:1`/`:0`).
  * `EXISTS`: Key presence verification with dynamic expiration handling.
  * `INCR`: Atomic numeric counter increments with robust type-checking and zero-initialization.

---

## Getting Started

### Prerequisites
* Go (version 1.18 or higher recommended)

### Running the Server
1. Clone the repository and navigate to the project directory.
2. Start the TCP server:
   ```bash
   go run main.go


3. Test the commands:
```redis
PING
ECHO hi there
SET greeting "hello"
GET greeting
SET countdown 5 EX 10
INCR countdown
EXISTS countdown
DEL greeting    