package binarypack

import (
	"errors"
	"fmt"
	"math"
)

// var BufferBuilder = require('./bufferbuilder').BufferBuilder;
// var binaryFeatures = require('./bufferbuilder').binaryFeatures;

const undefined = false

func Unpack(data []byte) interface{} {
	unpacker := newUnpacker(data)
	return unpacker.unpack()
}

func newUnpacker(data []byte) unpacker {
	return unpacker{
		index:      0,
		dataBuffer: data,
		// dataView: new Uint8Array(u.dataBuffer),
		dataView: data,
		length:   len(data),
	}
}

type unpacker struct {
	index      int
	dataBuffer []byte
	dataView   []byte
	length     int
}

func (u *unpacker) unpack() interface{} {
	ttype := u.unpack_uint8()
	if ttype < 0x80 {
		return ttype
	} else if (ttype ^ 0xe0) < 0x20 {
		return (ttype ^ 0xe0) - 0x20
	}

	size := ttype ^ 0xa0
	if size <= 0x0f {
		return u.unpack_raw(int(size))
	}

	size = ttype ^ 0xb0
	if size <= 0x0f {
		return u.unpack_string(int(size))
	}

	size = ttype ^ 0x90
	if size <= 0x0f {
		return u.unpack_array(int(size))
	}

	size = ttype ^ 0x80
	if size <= 0x0f {
		return u.unpack_map(int(size))
	}

	switch ttype {
	case 0xc0:
		return nil
	case 0xc1:
		return undefined
	case 0xc2:
		return false
	case 0xc3:
		return true
	case 0xca:
		return u.unpack_float()
	case 0xcb:
		return u.unpack_double()
	case 0xcc:
		return u.unpack_uint8()
	case 0xcd:
		return u.unpack_uint16()
	case 0xce:
		return u.unpack_uint32()
	case 0xcf:
		return u.unpack_uint64()
	case 0xd0:
		return u.unpack_int8()
	case 0xd1:
		return u.unpack_int16()
	case 0xd2:
		return u.unpack_int32()
	case 0xd3:
		return u.unpack_int64()
	case 0xd4:
		return undefined
	case 0xd5:
		return undefined
	case 0xd6:
		return undefined
	case 0xd7:
		return undefined
	case 0xd8:
		size := u.unpack_uint16()
		return u.unpack_string(int(size))
	case 0xd9:
		size := u.unpack_uint32()
		return u.unpack_string(int(size))
	case 0xda:
		size := u.unpack_uint16()
		return u.unpack_raw(int(size))
	case 0xdb:
		size := u.unpack_uint32()
		return u.unpack_raw(int(size))
	case 0xdc:
		size := u.unpack_uint16()
		return u.unpack_array(int(size))
	case 0xdd:
		size := u.unpack_uint32()
		return u.unpack_array(int(size))
	case 0xde:
		size := u.unpack_uint16()
		return u.unpack_map(int(size))
	case 0xdf:
		size := u.unpack_uint32()
		return u.unpack_map(int(size))
	}

	return undefined
}

func (u *unpacker) unpack_uint8() uint8 {
	b := u.dataView[u.index] & 0xff
	u.index++
	return b
}

func (u *unpacker) unpack_uint16() uint16 {
	bytes := u.read(2)

	b0 := uint16(bytes[0])
	b1 := uint16(bytes[1])

	uint16val := ((b0 & 0xff) * 256) + (b1 & 0xff)
	u.index += 2

	return uint16val
}

func (u *unpacker) unpack_uint32() uint32 {
	bytesRaw := u.read(4)
	bytes := make([]uint32, 8)
	for i, b := range bytesRaw {
		bytes[i] = uint32(b)
	}
	uint32val :=
		((bytes[0]*256+
			bytes[1])*256+
			bytes[2])*256 +
			bytes[3]
	u.index += 4
	return uint32val
}

func (u *unpacker) unpack_uint64() uint64 {
	bytesRaw := u.read(8)
	bytes := make([]uint64, 8)
	for i, b := range bytesRaw {
		bytes[i] = uint64(b)
	}
	uint64val :=
		((((((bytes[0]*256+
			bytes[1])*256+
			bytes[2])*256+
			bytes[3])*256+
			bytes[4])*256+
			bytes[5])*256+
			bytes[6])*256 +
			bytes[7]
	u.index += 8
	return uint64val
}

func (u *unpacker) unpack_int8() int8 {
	val := int(u.unpack_uint8())
	if val < 0x80 {
		return int8(val)
	}
	return int8(val - (1 << 8))
}

func (u *unpacker) unpack_int16() int16 {
	val := u.unpack_uint16()
	if val < 0x8000 {
		return int16(val)
	}
	return int16(int(val) - (1 << 16))
}

func (u *unpacker) unpack_int32() int32 {
	uint32val := u.unpack_uint32()
	if uint32val < math.Pow(2, 31) {
		return uint32val
	}
	return uint32val - math.Pow(2, 32)
}

func (u *unpacker) unpack_int64() int64 {
	uint64val := u.unpack_uint64()
	if uint64 < Math.pow(2, 63) {
		return uint64val
	}
	return uint64val - math.Pow(2, 64)
}

func (u *unpacker) unpack_raw(size int) []byte {
	if u.length < u.index+size {
		panic(fmt.Errorf(
			"BinaryPackFailure: index is out of range %d %d %d",
			u.index, size, u.length,
		))
	}
	buf := u.dataBuffer[u.index : u.index+size]
	u.index += size
	// buf = util.bufferToString(buf);
	return buf
}

func (u *unpacker) unpack_string(size int) string {
	bytes := u.read(int(size))
	i := 0
	str := ""
	var c byte
	var code rune

	for i < size {
		c = bytes[i]
		if c < 128 {
			//   str += String.fromCharCode(c);
			str += string(c)
			i++
		} else if (c ^ 0xc0) < 32 {
			code = rune((c^0xc0)<<6) | (bytes[i+1] & 63)
			//   str += String.fromCharCode(code);
			str += string(code)
			i += 2
		} else {
			code = ((c & 15) << 12) | ((bytes[i+1] & 63) << 6) |
				(bytes[i+2] & 63)
			str += string(code)
			i += 3
		}
	}

	u.index += size
	return str
}

func (u *unpacker) unpack_array(size int) []interface{} {
	objects := make([]interface{}, size)
	i := 0
	for i < size {
		objects[i] = u.unpack()
		i++
	}
	return objects
}

func (u *unpacker) unpack_map(size int) map[interface{}]interface{} {
	mapval := map[interface{}]interface{}{}
	i := 0
	for i < size {
		key := u.unpack()
		value := u.unpack()
		mapval[key] = value
		i++
	}
	return mapval
}

func (u *unpacker) unpack_float() float64 {
	uint32val := u.unpack_uint32()
	sign := uint32val >> 31
	exp := int(((uint32val >> 23) & 0xff) - 127)
	fraction := (uint32val & 0x7fffff) | 0x800000

	singval := -1
	if sign == 0 {
		singval = 1
	}

	return float64(singval) * float64(fraction) * math.Pow(2, float64(exp-23))
}

func (u *unpacker) unpack_double() float64 {
	h32 := u.unpack_uint32()
	l32 := u.unpack_uint32()
	sign := h32 >> 31
	exp := float64(((h32 >> 20) & 0x7ff) - 1023)
	hfrac := (h32 & 0xfffff) | 0x100000
	frac := float64(hfrac)*math.Pow(2, exp-20) + float64(l32)*math.Pow(2, exp-52)

	var signval float64 = -1
	if sign == 0 {
		signval = 1
	}

	return signval * frac
}

func (u *unpacker) read(length int) []byte {
	j := u.index
	if j+length <= u.length {
		return u.dataView[j : j+length]
	}
	panic(errors.New("BinaryPackFailure: read index out of range"))
}
