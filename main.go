package main

import (
	"log"
	"net"
)

const Port = "8000"
const StreamerMode = true

func safeRemoteAddr(conn net.Conn) {
	if StreamerMode {
		return "[REDACTED]"
	} else {
		return conn.RemoteAddr().String()
	}
}

func handleConnection(conn net.Conn) {
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
		go handleConnection(conn)
	}

}
