package connection

import (
	"bytes"
	"net"
	"t1/dispatch"
	"t1/logging"
)

type Client struct {
	ID   string
	Conn *net.TCPConn
	Pool *Pool
}

func (c *Client) Read() {
	logging.Println("Received connection from " + c.Conn.RemoteAddr().String())
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

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

	d_err := dispatch.Dispatch(c.Conn, packet)
	if d_err != nil {
		logging.Errorln("dispatch error:", d_err.Error())
		// logging.Infoln("remaining buffer:\n", hex.Dump(readBuf)) // end up being mostly nulls
		c.Conn.Close()
	}
}
