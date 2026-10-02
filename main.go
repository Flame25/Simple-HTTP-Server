package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	method  string
	path    string
	version string
	headers map[string]string
	body    string
}

type Status int

const (
	BadRequest Status = iota
	NotFound
	OK
	MethodNotAllowed
)

func writeResponse(conn net.Conn, body string, status Status) {

	statusText := ""
	switch status {
	case OK:
		statusText = "200 OK"
	case BadRequest:
		statusText = "400 Bad Request"
		body = "Bad Request"
	case NotFound:
		statusText = "404 Not Found"
		body = "Not Found"
	case MethodNotAllowed:
		statusText = "405 Method Not Allowed"
		body = "Method Not Allowed"
	}
	fmt.Fprintf(conn, "HTTP/1.1 %s\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%s", statusText, len(body), body)
}

func readRequest(reader *bufio.Reader) (Request, error) {

	currMsg, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return Request{}, err
	}

	requestMsg := strings.Split(currMsg, " ")
	if len(requestMsg) != 3 {
		fmt.Printf("Error: Header length not 3\n")
		return Request{}, fmt.Errorf("Invalid Request Line")
	} else if len(requestMsg[1]) <= 0 {
		fmt.Printf("Error: Path too short\n")
		return Request{}, fmt.Errorf("Invalid Path Line")
	}

	req := Request{
		version: strings.Trim(requestMsg[2], "\r\n"),
		method:  requestMsg[0],
		path:    requestMsg[1]}
	fmt.Printf("%q, %q, %q\n", req.method, req.path, req.version)

	return req, nil
}

func readBody(reader *bufio.Reader, num int) (string, error) {
	contentInput := make([]byte, num)
	n, err := io.ReadFull(reader, contentInput)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	b := string(contentInput)
	fmt.Printf("Len: %v, %s\n", n, b)
	return b, nil
}

func readHeader(reader *bufio.Reader) (map[string]string, error) {

	headers := make(map[string]string)
	var msg strings.Builder

	headerMsg := []string{}

	var err error
	currMsg := ""

	for currMsg != "\r\n" {
		currMsg, err = reader.ReadString('\n')
		msg.WriteString(currMsg)
		fmt.Printf("%q\n", currMsg)
		if currMsg != "\r\n" {
			headerMsg = append(headerMsg, currMsg)
		}

		if err != nil && err != io.EOF {
			return headers, err
		}

		if err == io.EOF {
			return headers, err
		}
	}

	fmt.Println("--- Content : --- ")
	fmt.Printf("%s", msg.String())
	fmt.Println()

	for _, el := range headerMsg {
		ns := strings.Trim(el, "\r\n")
		fmt.Printf("%q\n", ns)
		temp := strings.Split(ns, ": ")
		if len(temp) != 2 {
			fmt.Printf("Error: Invalid header\n")
			return headers, fmt.Errorf("Invalid Header")
		}
		headers[temp[0]] = temp[1]
		fmt.Printf("%s: %s\n", temp[0], temp[1])
	}

	// check if host exist needed for http
	_, ok := headers["Host"]
	if !ok {
		fmt.Printf("Error: Host not exist in header\n")
		return headers, fmt.Errorf("Host Invalid")
	}

	return headers, nil
}

func handleTcp(conn net.Conn) {
	fmt.Printf("Connected: %s\n", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		req, err := readRequest(reader)

		if err == io.EOF {
			return
		}

		if err != nil {
			log.Printf("Error: %s\n", err)
			writeResponse(conn, "", BadRequest)
			return
		}

		req.headers, err = readHeader(reader)

		if err == io.EOF {
			return
		}

		if err != nil {
			log.Printf("Error: %s\n", err)
			writeResponse(conn, "", BadRequest)
			return
		}

		switch req.method {
		case "GET":
			text := ""
			status := OK

			switch req.path {
			case "/":
				text = "Welcome!"
				status = OK
			case "/hello":
				text = "Hello!"
				status = OK
			case "/about":
				text = "About this server"
				status = OK
			default:
				text = "Not Found"
				status = NotFound
			}
			writeResponse(conn, text, status)
		case "POST":
			bodyLen, ok := req.headers["Content-Length"]
			if ok {
				temp, err := strconv.Atoi(bodyLen)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					writeResponse(conn, "", BadRequest)
					return
				} else if temp <= 0 {
					fmt.Printf("Error: Body too short\n")
					writeResponse(conn, "", BadRequest)
					return
				} else if temp > 10000 {
					fmt.Printf("Error: Body too long\n")
					writeResponse(conn, "", BadRequest)
					return
				}

				conn.SetReadDeadline(time.Now().Add(3 * time.Second))
				req.body, err = readBody(reader, temp)

				if err == io.EOF {
					return
				}

				if err != nil {
					fmt.Printf("Error: %s\n", err)
					writeResponse(conn, "", BadRequest)
					return
				}

				if req.path == "/echo" {
					writeResponse(conn, req.body, OK)
				} else {
					writeResponse(conn, "", NotFound)
				}
			} else {
				writeResponse(conn, "", BadRequest)
			}
		default:
			writeResponse(conn, "", MethodNotAllowed)
		}

		// check header for connection close
		conClose, ok := req.headers["Connection"]
		if ok {
			if conClose == "close" {
				return
			}
		}
	}
}

func main() {
	fmt.Println("--- TCP Server ---")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleTcp(conn)
	}
}
