package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"reflect"
	"t1/logging"
)

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
