package connection

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"t1/logging"
	"t1/packets"
)

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
