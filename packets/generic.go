package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"reflect"
	"t1/logging"
)

type BNCSBase struct {
	seek_index int `default:"0"`
	BNCSHeader
}

type BNCSHeader struct {
	Marker uint8 `default:"uint8(0xff)"`
	ID     uint8
	Length uint16
}

type BNCSGeneric struct {
	BNCSBase
	Data []byte
}

func (d BNCSBase) String() string {
	return fmt.Sprintf("%x:%x:%x -> %d", d.Marker, d.ID, d.Length, d.Length)
}

// Traverse a thing and get it back as a byte slice
// Technically works on basically anything, but mostly useful for converting packets between
// BNCSGeneric and a specific struct. Also for dumping to logs
func GetBytes(d interface{}) []byte {
	var err error
	buf := new(bytes.Buffer)

	dv := reflect.ValueOf(d)

	values := make([]interface{}, dv.NumField())
	for i := range values {
		if !dv.Field(i).CanInterface() {
			logging.Debugf("Field %d is unexported, skipping", i)
			continue
		}
		values[i] = dv.Field(i).Interface()

		switch dv.Field(i).Kind() {
		case reflect.Slice:
			slice_data := reflect.ValueOf(dv.Field(i).Interface()).Interface()
			err = binary.Write(buf, binary.LittleEndian, slice_data)

		case reflect.Struct:
			struct_data := GetBytes(dv.Field(i).Interface())
			err = binary.Write(buf, binary.LittleEndian, struct_data)

		case reflect.String:
			var buf2 bytes.Buffer
			enc := gob.NewEncoder(&buf2)
			err = enc.Encode(dv.Field(i).Interface())
			if err != nil {
				logging.Warningln("failed to encode string as bytes")
				continue
			}

			bs := buf2.Bytes()
			if bs[len(bs)-1] != 0x00 {
				// make sure the string is null terminated
				bs = append(bs, 0x00)
			}
			err = binary.Write(buf, binary.LittleEndian, bs)

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = binary.Write(buf, binary.LittleEndian, dv.Field(i).Interface())

		default:
			logging.Warningln("failed to convert field: ", dv.Type().Field(i).Name, dv.Field(i).Kind())
		}
		if err != nil {
			logging.Errorln("binary.Write failed:", err)
			err = nil
		}
	}

	return buf.Bytes()
}

func (d *BNCSGeneric) ResetSeek() {
	d.seek_index = 0
}

func (d *BNCSGeneric) SetSeek(position int) {
	d.seek_index = position
}

func (d *BNCSGeneric) GetSeek() int {
	return d.seek_index
}

func (d *BNCSGeneric) SetLength() {
	d.Length = uint16(4 + len(d.Data))
}

func (d *BNCSGeneric) WriteBytes(val []byte) {
	d.Data = append(d.Data, val...)
}

func (d *BNCSGeneric) ReadUint8() uint8 {
	val := uint8(d.Data[d.seek_index])
	d.seek_index += 1
	return val
}

func (d *BNCSGeneric) WriteUint8(val uint8) {
	d.Data = append(d.Data, val)
}

func (d *BNCSGeneric) ReadUint16() uint16 {
	val := uint16(binary.LittleEndian.Uint16(d.Data[d.seek_index : d.seek_index+2]))
	d.seek_index += 2
	return val
}

func (d *BNCSGeneric) WriteUint16(val uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, val)
	d.Data = append(d.Data, b...)
}

func (d *BNCSGeneric) ReadUint32() uint32 {
	val := uint32(binary.LittleEndian.Uint32(d.Data[d.seek_index : d.seek_index+4]))
	d.seek_index += 4
	return val
}

func (d *BNCSGeneric) WriteUint32(val uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, val)
	d.Data = append(d.Data, b...)
}

// BNCS is jank like this
func (d *BNCSGeneric) WriteUint32String(val string) {
	if len(val) > 4 {
		val = val[0:4]
	}
	b := make([]byte, 0)
	b = append(b, []byte(val)...)
	conv := binary.LittleEndian.Uint32(b)
	d.WriteUint32(conv)
}

// more jank for writing to general uint32s
func WriteUint32String(val string) uint32 {
	if len(val) > 4 {
		val = val[0:4]
	}
	b := make([]byte, 0)
	b = append(b, []byte(val)...)
	return binary.BigEndian.Uint32(b)
}

func (d *BNCSGeneric) ReadUint32String() string {
	raw := d.ReadUint32()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, raw)
	return string(b)
}

func (d *BNCSGeneric) ReadUint64() uint64 {
	val := uint64(binary.LittleEndian.Uint64(d.Data[d.seek_index : d.seek_index+8]))
	d.seek_index += 8
	return val
}

func (d *BNCSGeneric) WriteUint64(val uint64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)
	d.Data = append(d.Data, b...)
}

func (d *BNCSGeneric) ReadString() string {
	var ret string

	for i := range d.Data[d.seek_index:] {
		if d.Data[d.seek_index+i] == 0x00 {
			ret = string(d.Data[d.seek_index : d.seek_index+i])
			d.seek_index += i + 1
			return ret
		}
	}
	return ret
}

func (d *BNCSGeneric) WriteString(val string) {
	d.Data = append(d.Data, []byte(val)...)
	d.Data = append(d.Data, 0x00)
}
