package main

// TODO: make logging not suck

import (
	"log"
	"net"

	"t1/connection"
	"t1/logging"
)

// Main loop for handling incoming connections
func main() {
	logging.InitLogs()
	log.Println("Starting server")
	listener, _ := net.ListenTCP("tcp", &net.TCPAddr{Port: 9999})

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Println("Could not connect", err.Error())
			continue
		}

		go connection.ProcessorRoutine(conn)
	}
}
