package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

const LOGON_XSHA1 uint8 = 0x00
const LOGON_NLS1 uint8 = 0x01
const LOGON_NLS2 uint8 = 0x02

type BNCS_SID_AUTH_INFO struct {
	Marker          uint8
	ID              uint8
	Length          uint16
	Data            []byte
	LogonType       uint32
	ServerToken     uint32
	UDPValue        uint32
	CR_MPQ_Filetime uint64
	CR_MPQ_Filename string
	CR_Formula      string
	ServerSignature [128]byte // WAR3/W3XP only
}

func (d BNCS_SID_AUTH_INFO) String() string {
	return fmt.Sprintf("%x:%x:%x -> %d", d.Marker, d.ID, d.Length, d.Length)
}

// Convert a BNCSGeneric struct into a slice of bytes. Useful for sending packets
// as well as debugging with hex.Dump()
func (d BNCS_SID_AUTH_INFO) ToBytes() []byte {
	buf := new(bytes.Buffer)

	v := reflect.ValueOf(d)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	for _, v := range values {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}
	return buf.Bytes()
}

func (d BNCSGeneric) CoerceFrom() BNCS_SID_AUTH_INFO {
	auth_info := BNCS_SID_AUTH_INFO{}
	auth_info.Marker = d.Marker
	auth_info.ID = d.ID
	auth_info.Length = d.Length

	return auth_info
}
