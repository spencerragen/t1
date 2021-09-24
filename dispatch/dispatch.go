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

func validatePacketSize(conn *net.TCPConn, packet packets.BNCSGeneric) bool {
	data := packets.GetBytes(packet)
	if uint16(len(data)) != packet.Length {
		log.Println("[*] client", conn.RemoteAddr().String(), " - packet size mismatch: ", len(data), packet.Length)
		conn.Close()
		return false
	}

	return true
}

func Dispatch(conn *net.TCPConn, packet packets.BNCSGeneric) error {
	if !validatePacketSize(conn, packet) {
		log.Println("[*] client", conn.RemoteAddr().String(), "connection terminated")
		return nil
	}
	switch packet.ID {
	case SID_AUTH_INFO:
		log.Println("received SID_AUTH_INFO")
		p := packets.BNCS_CLIENT_SID_AUTH_INFO{}
		p.From(packet)

		response, err := p.Process()
		if err != nil {
			if strings.Contains(err.Error(), "configured unsupported") {
				packets.SendMessageBox(conn, "Client is unsupported", "Error")
			}
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
