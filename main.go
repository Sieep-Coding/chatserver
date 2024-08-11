package main

import (
	"log"
	"net"
)

const Port = "8000"
const SafeMode = true

func safeRemoteAddr(conn net.Conn) string {
	if SafeMode {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

func handleConnection(conn net.Conn, outgoing chan string) {
	defer conn.Close()
	message := []byte("Hello World\n")
	a, err := conn.Write(message)
	if err != nil {
		log.Printf("Could not write message to %s: %s\n", safeRemoteAddr(conn), err) // log address connecting to server
		return
	}
	if a < len(message) {
		log.Printf("The message was not fully written %d/%d\n", len(message))
		return
	}
	buffer := make([]byte, 512)
	for {
		a, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		outgoing <- string(buffer[0:a])
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s\n", Port, err)
	}
	log.Printf("Listening to TCP connections on port %s ...\n", Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("ERROR: Could not accept connection: %s\n", err)
			// handle error
		}
		log.Printf("Accepted Connection from %s\n", safeRemoteAddr(conn))
		outgoing := make(chan string)
		go handleConnection(conn, outgoing)
	}

}
