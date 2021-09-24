package connection

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"

	"t1/dispatch"
	"t1/logging"
	"t1/packets"
)

// Read data from a connection into a buffer and pass it into the packet builder
func ProcessorRoutine(conn *net.TCPConn) {
	logging.Infoln("Received connection from " + conn.RemoteAddr().String())
	defer conn.Close()

	localBuffer := new(bytes.Buffer)
	readBuf := make([]byte, 1024)
	dataLen, err := conn.Read(readBuf)
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

	d_err := dispatch.Dispatch(conn, packet)
	if d_err != nil {
		logging.Errorln("dispatch error:", d_err.Error())
		return
	}
}

// Take a byte buffer and coerce it into generic packet for processing
func ReadPacket(r *bytes.Buffer) (packets.BNCSGeneric, error) {
	check := make([]byte, 4)
	_, err := r.Read(check)
	if err != nil {
		return packets.BNCSGeneric{}, err
	}

	if check[0] != 0xff {
		return packets.BNCSGeneric{}, fmt.Errorf("sanity byte mismatch:\n%s", hex.Dump(check))
	}
	packetsize := int(binary.LittleEndian.Uint16([]byte{check[2], check[3]}))
	packetbuffer := make([]byte, packetsize-4)
	_, err = r.Read(packetbuffer)
	if err != nil {
		return packets.BNCSGeneric{}, err
	}

	ret := packets.BNCSGeneric{}
	ret.Marker = check[0]
	ret.ID = check[1]
	ret.Length = uint16(packetsize)
	ret.Data = packetbuffer
	logging.Infoln("Packet received:\n", hex.Dump(packets.GetBytes(ret)))

	return ret, nil
}
