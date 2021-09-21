package dispatch

import (
	"fmt"
	"log"
	"net"
	"strings"
	"t1/packets"
)

const (
	SID_AUTH_INFO uint8 = 0x50
)

func Dispatch(conn *net.TCPConn, packet packets.BNCSGeneric) error {
	switch packet.ID {
	case SID_AUTH_INFO:
		log.Println("received SID_AUTH_INFO")
		p := packets.BNCS_SERVER_SID_AUTH_INFO{}
		p.CoerceFrom(packet)
		response, err := p.Process()
		if err != nil {
			if strings.Contains(err.Error(), "configured unsupported") {
				packets.SendMessageBox(conn, "Client is unsupported", "Error")
			}
			conn.Close()
			return fmt.Errorf("invalid SID_AUTH_INFO: %v", err.Error())
		}
		send_response(conn, response)
	}
	return nil
}

func send_response(conn *net.TCPConn, response packets.BNCSGeneric) error {
	fmt.Println("send_response")
	return nil
}
