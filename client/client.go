package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Disconnected from server:", err)
			return
		}
		fmt.Print(string(buffer[:n]))
	}
}

func write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Error reasing input: ", err)
			return
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message: ", err)
			return
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	//TODO Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//TODO Start asynchronously reading and displaying messages
	go read(conn)
	//TODO Start getting and sending user messages.
	write(conn)
}
