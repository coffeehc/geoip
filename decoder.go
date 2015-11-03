// decoder
package geoip

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/coffeehc/logger"
)

const (
	TYPE_EXTENDED int = iota
	TYPE_POINTER
	TYPE_UTF8_STRING
	TYPE_DOUBLE
	TYPE_BYTES
	TYPE_UINT16
	TYPE_UINT32
	TYPE_MAP
	TYPE_INT32
	TYPE_UINT64
	TYPE_UINT128
	TYPE_ARRAY
	TYPE_CONTAINER
	TYPE_END_MARKER
	TYPE_BOOLEAN
	TYPE_FLOAT
)

var pointerValueOffset []int

func init() {
	pointerValueOffset = []int{0, 0, 1 << 11, (1 << 19) + (1 << 11), 0}
}

type decoder struct {
	pointerBase  int
	objectMapper map[string]interface{}
	database     *IpDataBase
	pointCache   map[int]*entry_data_cache
}

type entry_data_cache struct {
	node     interface{} //JsonNode
	nodeType int
}

type entry_data struct {
	node       interface{} //JsonNode
	offsetNext int
	nodeType   int
}

func (this *entry_data) getInt() int {
	if v, ok := this.node.(int); ok {
		return v
	}
	return 0
}

func newDecoder(database *IpDataBase, pointerBase int) *decoder {
	d := new(decoder)
	d.pointerBase = pointerBase
	d.database = database
	d.pointCache = make(map[int]*entry_data_cache)
	d.pointCache[-1] = nil
	return d
}

func (this *decoder) decode(offset int) (*entry_data, error) {
	var err error
	if offset >= len(this.database.data) {
		return nil, fmt.Errorf("解码边界超出数据长度,offSet:%d,数据长度:%d", offset, len(this.database.data))
	}
	ctrl := this.database.data[offset]
	offset++
	nodeType := int((ctrl >> 5) & 7)
	if nodeType == TYPE_EXTENDED {
		nodeType = get_ext_type(int(this.database.data[offset]))
		offset++
	}
	if nodeType == TYPE_POINTER {
		tmpEntry := this.decodePointer(ctrl, this.database.data[offset:], offset)
		if v, ok := tmpEntry.node.(int); ok {
			entry := new(entry_data)
			pc := this.pointCache[v]
			if pc != nil {
				entry.node = pc.node
				entry.nodeType = pc.nodeType
			} else {
				entry, err = this.decode(v)
				if err != nil {
					return nil, err
				}
				cache := new(entry_data_cache)
				cache.node = entry.node
				cache.nodeType = entry.nodeType
				this.pointCache[v] = cache
			}
			entry.offsetNext = tmpEntry.offsetNext
			return entry, nil
		} else {
			return nil, fmt.Errorf("解析Pointer异常")
		}
	}
	size := int(ctrl & 31)
	switch size {
	case 29:
		size = 29 + int(this.database.data[offset])
		offset++
		break
	case 30:
		size = 285 + int(get_uint16(this.database.data[offset:]))
		offset += 2
		break
	case 31:
		size = 65821 + int(get_uint24(this.database.data[offset:]))
		offset += 3
	default:
		break
	}
	entry := new(entry_data)
	entry.nodeType = nodeType
	switch nodeType {
	case TYPE_MAP:
		m := make(map[interface{}]interface{}) // 不一定是String
		for i := 0; i < size; i++ {
			key, err := this.decode(offset)
			if err != nil {
				logger.Error("解析出错Key:%s", err)
				return nil, err
			}
			value, err := this.decode(key.offsetNext)
			if err != nil {
				logger.Error("解析出错Value:%s", err)
				return nil, err
			}
			m[key.node] = value.node
			offset = value.offsetNext
			continue
		}
		entry.node = m
		entry.offsetNext = offset
		return entry, nil
	case TYPE_ARRAY:
		a := make([]interface{}, 0)
		for i := 0; i < size; i++ {
			value, err := this.decode(offset)
			if err != nil {
				logger.Error("解析出错Value:%s", err)
				return nil, err
			}
			a = append(a, value.node)
			offset = value.offsetNext
		}
		entry.node = a
		entry.offsetNext = offset
		return entry, nil
	case TYPE_BOOLEAN:
		if size != 0 {
			entry.node = true
		} else {
			entry.node = false
		}
		entry.offsetNext = offset
		return entry, nil
	case TYPE_UINT16:
		if size > 2 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		entry.node = uint16(get_uintX(this.database.data[offset:offset+size], 0, size))
		break
	case TYPE_UINT32:
		if size > 4 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		entry.node = uint32(get_uintX(this.database.data[offset:offset+size], 0, size))
		break
	case TYPE_INT32:
		if size > 4 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		entry.node = get_sintX(this.database.data[offset:offset+size], size)
		break
	case TYPE_UINT64:
		if size > 8 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		entry.node = get_uintX(this.database.data[offset:offset+size], 0, size)
		break
	case TYPE_UINT128:
		if size > 16 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		bytes := make([]byte, size)
		copy(bytes, this.database.data[offset:offset+size])
		entry.node = bytes
		break
	case TYPE_FLOAT:
		if size != 4 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		size = 4
		entry.node = get_float(this.database.data[offset : offset+size])
		break
	case TYPE_DOUBLE:
		if size != 8 {
			return nil, fmt.Errorf("无效的数据类型")
		}
		size = 8
		entry.node = get_double(this.database.data[offset : offset+size])
		break
	case TYPE_UTF8_STRING:
		if size == 0 {
			entry.node = ""
		} else {
			entry.node = fmt.Sprintf("%s", this.database.data[offset:offset+size])
		}
		break
	case TYPE_BYTES:
		b := make([]byte, size)
		copy(b, this.database.data[offset:offset+size])
		entry.node = b
		break
	}
	entry.offsetNext = offset + size
	return entry, nil
}

func get_ext_type(raw_ext_type int) int {
	return 7 + raw_ext_type
}

func (this *decoder) decodePointer(ctrl byte, data []byte, offset int) *entry_data {
	psize := int((ctrl>>3)&3 + 1)
	var base byte = 0
	if psize != 4 {
		base = ctrl & 0x7
	}
	new_offset := int(get_uintX(data, uint64(base), psize))
	entry := new(entry_data)
	entry.node = new_offset + this.pointerBase + pointerValueOffset[psize]
	entry.nodeType = TYPE_POINTER
	entry.offsetNext = psize + offset
	return entry
}

func get_double(p []byte) float64 {
	var f float64
	err := binary.Read(bytes.NewBuffer(p), binary.BigEndian, &f)
	if err != nil {
		logger.Error("解析float失败")
		return 0
	}
	return f
}

func get_float(p []byte) float32 {
	var f float32
	err := binary.Read(bytes.NewBuffer(p), binary.BigEndian, &f)
	if err != nil {
		logger.Error("解析float失败")
		return 0
	}
	return f
}

func get_uint16(p []byte) int16 {
	return int16(p[0])*256 + int16(p[1])
}

func get_uint24(p []byte) int32 {
	return int32(p[0])*65536 + int32(p[1])*256 + int32(p[2])
}

func get_uint32(p []byte) int32 {
	return int32(p[0])*16777216 + int32(p[1])*65536 + int32(p[2])*256 + int32(p[3])
}
func get_uintX(p []byte, value uint64, length int) uint64 {
	for i := 0; length > 0; i++ {
		value <<= 8
		value += uint64(p[i])
		length--
	}
	return value
}
func get_sintX(p []byte, length int) int32 {
	return int32(get_uintX(p, 0, length))
}

func typeToName(num int) string {
	switch num {
	case 0:
		return "extended"
	case 1:
		return "pointer"
	case 2:
		return "utf8_string"
	case 3:
		return "double"
	case 4:
		return "bytes"
	case 5:
		return "uint16"
	case 6:
		return "uint32"
	case 7:
		return "map"
	case 8:
		return "int32"
	case 9:
		return "uint64"
	case 10:
		return "uint128"
	case 11:
		return "array"
	case 12:
		return "container"
	case 13:
		return "end_marker"
	case 14:
		return "boolean"
	case 15:
		return "float"
	default:
		return "unknown type"
	}
}
