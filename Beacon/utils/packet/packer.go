package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

func PackBytes(b []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(len(b)))
	buf.Write(b)
	return buf.Bytes()
}

func PackArray(array []any) ([]byte, error) {
	var packData []byte

	for i := range array {
		switch v := array[i].(type) {

		case []byte:
			packData = append(packData, v...)
			break

		case string:
			size := make([]byte, 4)
			val := v
			if len(val) != 0 {
				if !strings.HasSuffix(val, "\x00") {
					val += "\x00"
				}
			}
			binary.BigEndian.PutUint32(size, uint32(len(val)))
			packData = append(packData, size...)
			packData = append(packData, []byte(val)...)
			break

		case int8:
			packData = append(packData, byte(v))
			break

		case int16:
			num := make([]byte, 2)
			binary.BigEndian.PutUint16(num, uint16(v))
			packData = append(packData, num...)
			break

		case int32:
			num := make([]byte, 4)
			binary.BigEndian.PutUint32(num, uint32(v))
			packData = append(packData, num...)
			break

		case int64:
			num := make([]byte, 8)
			binary.BigEndian.PutUint64(num, uint64(v))
			packData = append(packData, num...)
			break

		case bool:
			var bt byte = 0
			if v {
				bt = 1
			}
			packData = append(packData, bt)
			break

		default:
			return nil, errors.New("PackArray unknown type")
		}
	}
	return packData, nil
}

func MakeFinalPacket(taskId uint, commandId uint, data []byte) []byte {
	var (
		array    []interface{}
		err      error
		packData []byte
	)

	array = append(array, int32(taskId))
	array = append(array, int32(commandId))
	array = append(array, data)

	packData, err = PackArray(array)
	if err != nil {
		return nil
	}

	finalData := PackBytes(packData)

	return finalData
}
