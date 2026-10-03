package main

import (
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8000")

	if err != nil {
		log.Fatalf("failed to start listener %v", err)
	}
	defer listen.Close()
	log.Printf("server started on port %v", 8000)

	for {
		conn, err := listen.Accept()

		if err != nil {
			log.Printf("server failed to accept connection: %v", err)
			continue
		}

		go handler(conn)
	}
}
