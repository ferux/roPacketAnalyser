package packet

import (
	"bytes"
	"encoding/binary"
)

//BytesToInt32 binary reads from sequence of bytes to return int32
func BytesToInt32(buf *bytes.Buffer) int32 {
	var i int32
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToInt16 binary reads from sequence of bytes to return int16
func BytesToInt16(buf *bytes.Buffer) int16 {
	var i int16
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToInt8 binary reads from sequence of bytes to return int8
func BytesToInt8(buf *bytes.Buffer) int8 {
	var i int8
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToUint32 binary reads from sequence of bytes to return uint32
func BytesToUint32(buf *bytes.Buffer) uint32 {
	var i uint32
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToUint16 binary reads from sequence of bytes to return int16
func BytesToUint16(buf *bytes.Buffer) uint16 {
	var i uint16
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToUint8 binary reads from sequence of bytes to return int8
func BytesToUint8(buf *bytes.Buffer) uint8 {
	var i uint8
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		panic(err)
	}
	return i
}

//BytesToByteArray creates a new byte array with specified length
func BytesToByteArray(buf *bytes.Buffer, length int) []byte {
	//this+0x6 unsigned char ID[24];
	i := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, i); err != nil {
		panic(err)
	}
	return i
}

/*
	p[0] p[1] p[2].
	x = p[0] << 2 = (0011 1100)
	x = p[1] & 0xC0 >> 6 | x
	y = p[1] << & 0x3FF | (p[2] &0xF0 >> 4)
	dir = p[2] & 0xf
*/
//BytesToXYDir converts byte array with length = 3 to x, y, dir
func BytesToXYDir(bytes []uint8) (x, y int16, d uint8) {
	if len(bytes) != 3 {
		return
	}
	x = (int16(bytes[0]) << 2) | (int16(bytes[1]) & 0xC0 >> 6)
	y = (int16(bytes[1]) << 4 & 0x3ff) | (int16(bytes[2]&0xF0) >> 4)
	d = bytes[2] & 0xF
	return
}
