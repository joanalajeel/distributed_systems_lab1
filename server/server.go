package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message, disconnected chan int) {
	reader := bufio.NewReader(client)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
			disconnected <- clientid
			client.Close()
			return
		}

		msg := Message{
			sender:  clientid,
			message: message,
		}
		msgs <- msg
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		panic(err)
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	// add disconnected
	disconnected := make(chan int)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	nextClientID := 1

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			clientid := nextClientID
			nextClientID++
			clients[clientid] = conn

			fmt.Println("Client connection:", clientid)

			go handleClient(conn, clientid, msgs, disconnected)

		case msg := <-msgs:
			for id, conn := range clients {
				if id != msg.sender {
					_, err := conn.Write([]byte(msg.message))
					if err != nil {
						handleError(err)
						conn.Close()
						delete(clients, id)
					}
				}
			}

		case id := <-disconnected:
			if conn, ok := clients[id]; ok {
				conn.Close()
				delete(clients, id)
				fmt.Println("Client disconnected:", id)
			}
		}
	}
}
