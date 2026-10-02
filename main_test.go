package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func Test_readRequest(t *testing.T) {
	input := "GET /hello HTTP/1.1\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	req, err := readRequest(reader)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.method != "GET" {
		t.Errorf("expected method GET, got %q", req.method)
	}

	if req.path != "/hello" {
		t.Errorf("expected path /hello, got %q", req.path)
	}

	if req.version != "HTTP/1.1" {
		t.Errorf("expected version HTTP/1.1, got %q", req.version)
	}
}

func Test_invalidRequest(t *testing.T) {
	input := "GET /hello\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	_, err := readRequest(reader)
	if err == nil {
		t.Errorf("expected error invalid request")
	}

	input = "GET HTTP/1.1\r\n"
	reader = bufio.NewReader(strings.NewReader(input))

	_, err = readRequest(reader)
	if err == nil {
		t.Errorf("expected error invalid request")
	}
}

func Test_readHeader(t *testing.T) {

	input := "Host: localhost:8080\r\nUser-Agent: curl/8.18.0\r\nAccept: */*\r\nContent-Length: 24\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	red, err := readHeader(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h, ok := red["Host"]
	if h != "localhost:8080" || !ok {
		t.Errorf("expected localhost:8080")
	}

	h, ok = red["Content-Length"]
	if h != "24" || !ok {
		t.Errorf("expected 24")
	}

}

func Test_invalidHeader(t *testing.T) {

	input := "Host localhost:8080\r\nUser-Agent: curl/8.18.0\r\nAccept: */*\r\nContent-Length: 24\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	_, err := readHeader(reader)

	if err == nil {
		t.Errorf("expected error invalid header")
	}

	input = "Hostss: localhost:8080\r\nUser-Agent: curl/8.18.0\r\nAccept: */*\r\nContent-Length: 24\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
	reader = bufio.NewReader(strings.NewReader(input))

	_, err = readHeader(reader)

	if err == nil {
		t.Errorf("expected error invalid host")
	}
}

func Test_readBody(t *testing.T) {

	input := []byte("Hello")
	reader := bufio.NewReader(bytes.NewReader(input))

	b, err := readBody(reader, len(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b != "Hello" {
		t.Errorf("expected output Hello")
	}

}
