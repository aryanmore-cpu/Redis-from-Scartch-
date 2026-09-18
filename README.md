
# Go-Redis: A Custom In-Memory Key-Value Store

A lightweight, Redis-compatible TCP server built from scratch in Go. This project explores TCP networking, concurrent client handling, custom protocol parsing, thread-safe data storage, and key expiration without external framework dependencies.

The goal is to understand how an in-memory database and Redis-style server work internally by implementing its core components from the ground up.

---

## Features Implemented

### Raw TCP Server & Concurrency

- Built using Go's `net` package and TCP sockets.
- Supports multiple concurrent client connections using goroutines.
- Handles client requests through dedicated connection handlers.

### Thread-Safe Architecture

- Uses `sync.RWMutex` to protect shared in-memory data.
- Supports concurrent access to the key-value store.
- Separates server logic from storage operations.

### Key-Value Storage & Expiration

- Stores key-value pairs in an in-memory Go map.
- Uses custom `Item` structures to track values and expiration timestamps.
- Supports time-to-live (TTL) expiration.
- Implements lazy expiration and cleanup when keys are accessed.
- Supports dynamically setting expiration on existing keys using `EXPIRE`.

### RESP Protocol Parser

- Supports Redis Serialization Protocol (RESP) array commands.
- Parses length-prefixed bulk strings.
- Handles commands such as:

  ```text
  *1\r\n$4\r\nPING\r\n
  ```

- Maintains compatibility with plain-text, space-separated commands for manual testing.

---

## Project Architecture

The project is organized into three main Go files:

```text
.
├── main.go
├── server.go
├── store.go
├── go.mod
└── README.md
```

### `main.go`

Responsible for:

- Application entry point.
- Starting the TCP listener.
- Bootstrapping the Redis-compatible server.

### `server.go`

Responsible for:

- Accepting and handling client connections.
- Managing client handler goroutines.
- Parsing incoming commands.
- Decoding RESP array requests.
- Routing commands to the appropriate operations.
- Generating RESP-compatible responses.

### `store.go`

Responsible for:

- Managing the in-memory key-value store.
- Defining the `Item` data structure.
- Handling thread-safe data access.
- Managing TTL timestamps.
- Performing lazy expiration and cleanup.
- Implementing storage-related operations.

---

## Supported Commands

| Command | Description |
|---|---|
| `PING` | Tests whether the server is responsive |
| `ECHO` | Returns a message to the client |
| `SET` | Stores a key-value pair |
| `GET` | Retrieves the value of a key |
| `SET ... EX` | Stores a key with a time-to-live |
| `EXPIRE` | Sets or updates the expiration time of an existing key |
| `DEL` | Deletes a key |
| `EXISTS` | Checks whether a key exists |
| `INCR` | Increments a numeric value atomically |

### Example Commands

```text
PING

ECHO hello

SET greeting hello

GET greeting

SET countdown 5 EX 10

EXPIRE countdown 20

EXISTS countdown

INCR countdown

DEL greeting
```

---

## RESP Protocol Support

The server supports Redis-style RESP array requests.

Example:

```text
*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n
```

The parser processes:

- Array headers (`*`)
- Bulk string lengths (`$`)
- Carriage return and line feed delimiters (`\r\n`)
- Command arguments and values

The server also supports plain-text commands for easier manual testing.

---

## Getting Started

### Prerequisites

- Go 1.18 or higher recommended
- Windows, Linux, or macOS
- A TCP client or Redis-compatible client for testing

### Clone the Repository

```bash
git clone https://github.com/aryanmore-cpu/Redis-from-Scartch-.git
```

Navigate to the project directory:

```bash
cd Redis-from-Scartch-
```

### Run the Server

Start the server using:

```bash
go run .
```

The server listens on TCP port `6379`.

---

## Testing

The project has been tested using TCP requests and PowerShell-based testing workflows.

Testing includes:

- Basic connectivity using `PING`
- Message responses using `ECHO`
- Key storage and retrieval using `SET` and `GET`
- TTL expiration using `SET ... EX`
- Dynamic expiration using `EXPIRE`
- Key deletion using `DEL`
- Key existence checks using `EXISTS`
- Numeric increments using `INCR`
- RESP array command requests
- Plain-text command compatibility

Example TCP connection:

```powershell
$client = New-Object System.Net.Sockets.TcpClient("localhost", 6379)
```

---

## Technical Concepts Explored

This project provides practical experience with:

- TCP networking
- Socket-based communication
- Go goroutines
- Concurrent client handling
- Mutexes and race-condition prevention
- In-memory data structures
- TTL and lazy expiration
- RESP serialization and parsing
- Command routing
- Git and GitHub-based development


