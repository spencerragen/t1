package connection

import (
	"fmt"
	"strings"
	"t1/logging"
	"t1/packets"
	"t1/utils"
)

const (
	SID_AUTH_INFO uint8 = 0x50
)

func validatePacketSize(c *Client, packet packets.BNCSGeneric) bool {
	data := utils.GetBytes(packet)
	if uint16(len(data)) != packet.Length {
		logging.Errorln("client", c.Conn.RemoteAddr().String(), " - packet size mismatch: ", len(data), packet.Length)
		c.Terminate()
		return false
	}

	return true
}

func Dispatch(c *Client, packet packets.BNCSGeneric) error {
	if !validatePacketSize(c, packet) {
		logging.Errorln("client", c.Conn.RemoteAddr().String(), "connection terminated")
		return nil
	}
	switch packet.ID {
	case SID_AUTH_INFO:
		logging.Infoln("received SID_AUTH_INFO")
		p := packets.BNCS_CLIENT_SID_AUTH_INFO{}
		p.From(packet)

		response, err := p.Process(&c.LocalIP)
		if err != nil {
			if strings.Contains(err.Error(), "configured unsupported") {
				packets.SendMessageBox(c.Conn, "Client is unsupported", "Error")
			}
			return fmt.Errorf("invalid SID_AUTH_INFO: %v", err.Error())
		}
		send_response(c, response)
	}
	return nil
}

func send_response(conn *Client, response packets.BNCSGeneric) error {
	fmt.Println("send_response")
	return nil
}
