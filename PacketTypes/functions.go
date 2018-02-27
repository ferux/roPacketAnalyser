package PacketTypes

import (
	"bytes"
	"encoding/binary"
)

//BytesToInt32 binary reads from sequence of bytes to return int32
func BytesToInt32(buf *bytes.Buffer) int32 {
	var i int32
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToInt16 binary reads from sequence of bytes to return int16
func BytesToInt16(buf *bytes.Buffer) int16 {
	var i int16
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToInt8 binary reads from sequence of bytes to return int8
func BytesToInt8(buf *bytes.Buffer) int8 {
	var i int8
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToUint32 binary reads from sequence of bytes to return uint32
func BytesToUint32(buf *bytes.Buffer) uint32 {
	var i uint32
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToUint16 binary reads from sequence of bytes to return int16
func BytesToUint16(buf *bytes.Buffer) uint16 {
	var i uint16
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToUint8 binary reads from sequence of bytes to return int8
func BytesToUint8(buf *bytes.Buffer) uint8 {
	var i uint8
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

//BytesToByteArray creates a new byte array with specified length
func BytesToByteArray(buf *bytes.Buffer, length int) []byte {
	//this+0x6 unsigned char ID[24];
	i := make([]byte, length)
	binary.Read(buf, binary.LittleEndian, i)
	return i
}
