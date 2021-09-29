package connection

import (
	"bytes"
	"net"
	"t1/logging"
)

type Client struct {
	ID          string
	Conn        *net.TCPConn
	Pool        *Pool
	Channel     string
	LocalIP     string
	ServerToken uint32
	// TZBias      int32 // signed? unsigned? who knows!
}

func (c *Client) Terminate() {
	c.Pool.Unregister <- c
	c.Conn.Close()
}

func (c *Client) Read() {
	logging.Println("Received connection from " + c.Conn.RemoteAddr().String())
	defer c.Terminate()

	localBuffer := new(bytes.Buffer)
	readBuf := make([]byte, 1024)
	dataLen, err := c.Conn.Read(readBuf)
	if err != nil {
		logging.Errorln("error reading from stream:", err.Error())
		return
	}

	// this is bad, as protocol packets can be split over many tcp packets, and protocol
	// packets can be crammed together in single tcp packets.
	localBuffer.Write(readBuf[:dataLen])
	packet, err := ReadPacket(localBuffer)
	if err != nil {
		logging.Errorln("error forming packet:", err.Error())
		return
	}

	d_err := Dispatch(c, packet)
	if d_err != nil {
		logging.Errorln("dispatch error:", d_err.Error())
		// logging.Infoln("remaining buffer:\n", hex.Dump(readBuf)) // end up being mostly nulls
		c.Terminate()
	}
}
