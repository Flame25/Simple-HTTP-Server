package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	method  string
	path    string
	version string
	headers map[string]string
	body    string
}

func handleTcp(conn net.Conn) {
	fmt.Printf("Connected: %s\n", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	currMsg, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	var msg strings.Builder
	req := Request{
		headers: make(map[string]string),
	}
	requestMsg := strings.Split(currMsg, " ")
	if len(requestMsg) != 3 {
		fmt.Printf("Error: Header length not 3\n")
		return
	}

	req.version = strings.Trim(requestMsg[2], "\r\n")
	req.method = requestMsg[0]
	req.path = requestMsg[1]
	fmt.Printf("%q, %q, %q\n", req.method, req.path, req.version)
	headerMsg := []string{}
	for currMsg != "\r\n" {
		msg.WriteString(currMsg)
		currMsg, err = reader.ReadString('\n')
		if currMsg != "\r\n" {
			headerMsg = append(headerMsg, currMsg)
		}

		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			return
		}
	}

	fmt.Println("--- Content : --- ")
	fmt.Printf("%s", msg.String())
	fmt.Println()

	for _, el := range headerMsg {
		ns := strings.Trim(el, "\r\n")
		fmt.Printf("%q\n", ns)
		temp := strings.SplitN(ns, ": ", 2)
		req.headers[temp[0]] = temp[1]
	}

	switch req.method {
	case "GET":
		text := ""
		status := ""

		switch req.path {
		case "/":
			text = "Welcome!"
			status = "200 OK"
		case "/hello":
			text = "Hello!"
			status = "200 OK"
		case "/about":
			text = "About this server"
			status = "200 OK"
		default:
			text = "Not Found"
			status = "404 Not Found"
		}
		fmt.Fprintf(conn, "HTTP/1.1 %s\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%s", status, len(text), text)
	case "POST":
		bodyLen, ok := req.headers["Content-Length"]
		if ok {
			temp, err := strconv.Atoi(bodyLen)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request")
				return
			}
			if req.path == "/echo" {
				contentInput := make([]byte, temp)
				n, err := io.ReadFull(reader, contentInput)
				if err != nil {
					fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request")
					fmt.Printf("Error: %s\n", err.Error())
					return
				}
				req.body = string(contentInput)
				fmt.Printf("Len: %v, %s\n", n, req.body)
				fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%s", len(req.body), req.body)

			} else {
				fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\nNot Found")
			}
		} else {
			fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nBad Request")
		}
	default:
		fmt.Fprintf(conn, "HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\n\r\nNot Allowed")
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
