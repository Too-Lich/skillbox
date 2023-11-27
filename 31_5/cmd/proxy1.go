package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8084")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	request := bufio.NewReader(conn)

	fmt.Println("New connection from", conn.RemoteAddr())

	url := "http://localhost:8081"

	if _, err := fmt.Fscan(request, &url); err != nil || url == "" {
		fmt.Printf("Invalid URL: %v\n", err)
		return
	} else if _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		return
	}
	// Реализуйте логику обработки запроса

	response, err := http.Get(url)
	if err != nil {
		fmt.Print("Error making request:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print("Error reading response body:", err)
		return
	}

	_, err = conn.Write(body)
	if err != nil {
		fmt.Print("Error writing to client:", err)
		return
	}
}
