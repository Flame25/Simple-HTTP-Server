# Go HTTP Server From Scratch

A minimal HTTP server built from scratch in Go using raw TCP connections, without using Go's `net/http` package.

This project is part of my journey to understand how HTTP servers work underneath the abstractions provided by standard libraries.

## Features

* TCP server using `net`
* HTTP request-line parsing
* HTTP header parsing
* Request body handling with `Content-Length`
* Basic routing
* GET and POST methods
* HTTP status codes

  * `200 OK`
  * `400 Bad Request`
  * `404 Not Found`
  * `405 Method Not Allowed`
* Multiple requests over the same TCP connection
* `Connection: close` support
* Read deadlines for request and body timeouts
* Basic HTML response
* Concurrent client handling with goroutines
* Parser tests and race detection

## Routes

| Method | Path     | Description                |
| ------ | -------- | -------------------------- |
| GET    | `/`      | Simple HTML page           |
| GET    | `/about` | Another Simple HTML page   |
| POST   | `/echo`  | Echoes the request body    |

## Running

Start the server:

```bash
go run .
```

The server listens on:

```text
localhost:8080
```

Open the following in a browser:

```text
http://localhost:8080/
```

Or test it with `curl`:

```bash
curl http://localhost:8080/hello
```

POST example:

```bash
curl -X POST http://localhost:8080/echo -d "Hello from curl"
```

## What I Learned

Building this server helped me understand how TCP streams data, how HTTP requests and connections are structured, how headers and request bodies are handled, and how buffering, timeouts.

## Testing

The request parser has tests covering valid and malformed:

* Request lines
* Headers
* Request bodies

The project can also be checked with Go's race detector:

```bash
go test -race ./...
```

## Purpose

This is a learning project rather than a production-ready HTTP implementation.

The goal is to understand the fundamentals of Go as well as HTTP, TCP networking, request parsing, connection management, and concurrent server design by implementing the pieces myself.
