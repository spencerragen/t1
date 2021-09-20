package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"t1/dispatch"
	"t1/packets"
)

func main() {
	fmt.Println("Starting server")
	listener, _ := net.ListenTCP("tcp", &net.TCPAddr{Port: 9999})

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Println("Could not connect", err.Error())
			continue
		}

		go processorRoutine(conn)
	}
}

func processorRoutine(conn *net.TCPConn) {
	fmt.Println("Received connection from " + conn.RemoteAddr().String())
	defer conn.Close()

	localBuffer := new(bytes.Buffer)
	readBuf := make([]byte, 1024)
	dataLen, err := conn.Read(readBuf)
	if err != nil {
		fmt.Println("[!] error reading from stream:", err.Error())
	}
	localBuffer.Write(readBuf[:dataLen])
	packet, err := ReadPacket(localBuffer)
	if err != nil {
		fmt.Println("[!] error forming packet:", err.Error())
	}

	d_err := dispatch.Dispatch(packet)
	if d_err != nil {
		fmt.Println("[!] dispatch error:", d_err.Error())
	}
}

func ReadPacket(r *bytes.Buffer) (packets.BNCSGeneric, error) {
	check := make([]byte, 4)
	_, err := r.Read(check)
	if err != nil {
		return packets.BNCSGeneric{}, err
	}
	fmt.Println(hex.Dump(check))

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
	return ret, nil
}
