package main

// TODO: make logging not suck

import (
	"net"

	"t1/connection"
	"t1/logging"
)

// Main loop for handling incoming connections
func main() {
	logging.InitLogs()
	logging.Infoln("Starting server")
	pool := connection.NewPool()
	go pool.Start()
	listener, _ := net.ListenTCP("tcp", &net.TCPAddr{Port: 9999})

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logging.Errorln("Could not connect", err.Error())
			continue
		}
		client := &connection.Client{
			Conn: conn,
			Pool: pool,
		}

		client.Read()
	}
}
