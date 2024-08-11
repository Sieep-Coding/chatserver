package main

import (
	"fmt"
	"log"
	"net"
)

const SafeMode = true
const Port = "8000"

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
	DeleteClient
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

func server(messages chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-messages
		switch msg.Type {
		case ClientConnected:
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
		case DeleteClient:
			delete(conns, msg.Conn.RemoteAddr().String())
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

func client(conn net.Conn, messages chan Message) {
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
			messages <- Message{
				Type: DeleteClient,
				Conn: conn,
			}
			return
		}
		messages <- Message{
			Type: NewMessage,
			Text: string(buffer[0:a]),
			Conn: conn,
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s\n", Port, err)
	}
	log.Printf("Listening to TCP connections on port %s ...\n", Port)
	messages := make(chan Message)
	go server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("ERROR: Could not accept connection: %s\n", err)
		}
		fmt.Printf("Accepted Connection from %s\n", safeRemoteAddr(conn))
		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		go client(conn, messages)
	}
}
