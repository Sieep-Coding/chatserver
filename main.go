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

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

func server(messages chan Message) {
	conns := []net.Conn{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			conns = append(conns, msg.Conn)
		case NewMessage:
			for _, conn := range conns {
				_, err := conn.Write([]byte(msg.Text))
				if err != nil {
					//TODO: Remove connection from list
					log.Println("Could not send data to %s: %s", safeRemoteAddr(conn), err)
				}
			}
		}
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
		}
		log.Printf("Accepted Connection from %s\n", safeRemoteAddr(conn))
		outgoing := make(chan string)
		go handleConnection(conn, outgoing)
	}

}
