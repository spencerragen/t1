package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type BNCSGeneric struct {
	Marker uint8
	ID     uint8
	Length uint16
	Data   []byte
}

var index int = 0

func (d BNCSGeneric) String() string {
	return fmt.Sprintf("%x:%x:%x -> %d", d.Marker, d.ID, d.Length, d.Length)
}

// Convert a BNCSGeneric struct into a slice of bytes. Useful for sending packets
// as well as debugging with hex.Dump()
func (d BNCSGeneric) ToBytes() []byte {
	buf := new(bytes.Buffer)

	v := reflect.ValueOf(d)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fmt.Println("Extracting:", v.Field(i).Type().Name())
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

func (d *BNCSGeneric) SetLength() {
	d.Length = uint16(4 + len(d.Data))
}

func (d *BNCSGeneric) WriteBytes(val []byte) {
	d.Data = append(d.Data, val...)
}

func (d *BNCSGeneric) ReadUint8() uint8 {
	val := uint8(d.Data[index])
	index += 1
	return val
}

func (d *BNCSGeneric) WriteUint8(val uint8) {
	d.Data = append(d.Data, val)
}

func (d *BNCSGeneric) ReadUint16() uint16 {
	val := uint16(binary.LittleEndian.Uint16(d.Data[index : index+2]))
	index += 2
	return val
}

func (d *BNCSGeneric) WriteUint16(val uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, val)
	d.Data = append(d.Data, b...)
}

func (d *BNCSGeneric) ReadUint32() uint32 {
	val := uint32(binary.LittleEndian.Uint32(d.Data[index : index+4]))
	index += 4
	return val
}

func (d *BNCSGeneric) WriteUint32(val uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, val)
	d.Data = append(d.Data, b...)
}

func (d *BNCSGeneric) ReadUint64() uint64 {
	val := uint64(binary.LittleEndian.Uint64(d.Data[index : index+8]))
	index += 8
	return val
}

func (d *BNCSGeneric) WriteUint64(val uint64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)
	d.Data = append(d.Data, b...)
}

func (d *BNCSGeneric) ReadString() string {
	var ret string

	for i := range d.Data[index:] {
		if d.Data[index+i] == 0x00 {
			ret = string(d.Data[index : index+i])
			index += i + 1
			return ret
		}
	}
	return ret
}

func (d *BNCSGeneric) WriteString(val string) {
	d.Data = append(d.Data, []byte(val)...)
	d.Data = append(d.Data, 0x00)
}
