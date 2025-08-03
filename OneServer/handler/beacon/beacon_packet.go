package beacon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type Parser struct {
	buffer []byte
}

func CreateParser(buffer []byte) *Parser {
	return &Parser{
		buffer: buffer,
	}
}

func (p *Parser) Size() uint {
	return uint(len(p.buffer))
}

func (p *Parser) Check(types []string) bool {

	packerSize := p.Size()

	for _, t := range types {
		switch t {

		case "byte":
			if packerSize < 1 {
				return false
			}
			packerSize -= 1

		case "int16":
			if packerSize < 2 {
				return false
			}
			packerSize -= 2

		case "int32":
			if packerSize < 4 {
				return false
			}
			packerSize -= 4

		case "int64":
			if packerSize < 8 {
				return false
			}
			packerSize -= 8

		case "array":
			if packerSize < 4 {
				return false
			}

			index := p.Size() - packerSize
			value := make([]byte, 4)
			copy(value, p.buffer[index:index+4])
			length := uint(binary.BigEndian.Uint32(value))
			packerSize -= 4

			if packerSize < length {
				return false
			}
			packerSize -= length
		}
	}
	return true
}

func (p *Parser) ParseInt8() uint8 {
	var value = make([]byte, 1)

	if p.Size() >= 1 {
		if p.Size() == 1 {
			copy(value, p.buffer[:p.Size()])
			p.buffer = []byte{}
		} else {
			copy(value, p.buffer[:1])
			p.buffer = p.buffer[1:]
		}
	} else {
		return 0
	}

	return value[0]
}

func (p *Parser) ParseInt16() uint16 {
	var value = make([]byte, 2)

	if p.Size() >= 2 {
		if p.Size() == 2 {
			copy(value, p.buffer[:p.Size()])
			p.buffer = []byte{}
		} else {
			copy(value, p.buffer[:2])
			p.buffer = p.buffer[2:]
		}
	} else {
		return 0
	}

	return binary.BigEndian.Uint16(value)
}

func (p *Parser) ParseInt32() uint {
	var value = make([]byte, 4)

	if p.Size() >= 4 {
		if p.Size() == 4 {
			copy(value, p.buffer[:p.Size()])
			p.buffer = []byte{}
		} else {
			copy(value, p.buffer[:4])
			p.buffer = p.buffer[4:]
		}
	} else {
		return 0
	}

	return uint(binary.BigEndian.Uint32(value))
}

func (p *Parser) ParseInt64() uint64 {
	var value = make([]byte, 8)

	if p.Size() >= 8 {
		if p.Size() == 8 {
			copy(value, p.buffer[:p.Size()])
			p.buffer = []byte{}
		} else {
			copy(value, p.buffer[:8])
			p.buffer = p.buffer[8:]
		}
	} else {
		return 0
	}

	return binary.BigEndian.Uint64(value)
}

func (p *Parser) ParseBytes() []byte {
	size := p.ParseInt32()

	if p.Size() < size {
		return make([]byte, 0)
	} else {
		b := p.buffer[:size]
		p.buffer = p.buffer[size:]
		return b
	}
}

func (p *Parser) ParseString() string {
	size := p.ParseInt32()

	if p.Size() < size {
		return ""
	} else {
		b := p.buffer[:size]
		p.buffer = p.buffer[size:]
		return string(bytes.Trim(b, "\x00"))
	}
}

func (p *Parser) ParseString64() string {
	size := p.ParseInt64() // 读取8字节长度

	if p.Size() < uint(size) {
		return ""
	} else {
		b := p.buffer[:size]
		p.buffer = p.buffer[size:]
		return string(bytes.Trim(b, "\x00"))
	}
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
