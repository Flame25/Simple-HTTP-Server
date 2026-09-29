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

func handleTcp(conn net.Conn) {
	fmt.Printf("Connected: %s\n", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	redis := make(map[string]string)
	for {
		currMsg, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		var msg strings.Builder
		requestMsg := strings.SplitN(currMsg, " ", 3)
		requestMsg[2] = strings.Trim(requestMsg[2], "\r\n")
		fmt.Printf("%q, %q, %q\n", requestMsg[0], requestMsg[1], requestMsg[2])
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
			redis[temp[0]] = temp[1]
		}

		fmt.Printf("len : %s\n", redis["Content-Length"])

		bodyLen, ok := redis["Content-Length"]
		temp, _ := strconv.Atoi(bodyLen)
		if ok {
			contentInput := make([]byte, temp)
			n, err := io.ReadFull(reader, contentInput)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Len: %v, %s\n", n, string(contentInput))
		}

		text := ""
		status := ""

		switch requestMsg[1] {
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
