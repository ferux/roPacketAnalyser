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

/*
	p[0] p[1] p[2].
	x = p[0] << 2 = (0011 1100)
	x = p[1] & 0xC0 >> 6 | x
	y = p[1] << & 0x3FF | (p[2] &0xF0 >> 4)
	dir = p[2] & 0xf
*/
//BytesToXYDir converts byte array with length = 3 to x, y, dir
func BytesToXYDir(bytes []uint8) (x, y int16, d uint8) {
	if len(bytes) == 3 { //int16 x0, y0, direction
		x = (int16(bytes[0]) << 2) | (int16(bytes[1]) & 0xC0 >> 6)
		y = (int16(bytes[1]) << 4 & 0x3ff) | (int16(bytes[2]&0xF0) >> 4)
		d = bytes[2] & 0xF
	}
	if len(bytes) == 6 { //int16 x0, y0, x1, y1, uint8 sx0, sy0
		x = (int16(bytes[0]) << 2) | (int16(bytes[1]) & 0xC0 >> 6)
		y = (int16(bytes[1]) << 4 & 0x3ff) | (int16(bytes[2]&0xF0) >> 4)
		d = bytes[2] & 0xF
		//TODO
	}
	return
}

/*
	p[0] = (uint8)(x0>>2);
	p[1] = (uint8)((x0<<6) | ((y0>>4)&0x3f));
	p[2] = (uint8)((y0<<4) | ((x1>>6)&0x0f));

	p[3] = (uint8)((x1<<2) | ((y1>>8)&0x03));
	p[4] = (uint8)y1;
	p[5] = (uint8)((sx0<<4) | (sy0&0x0f));

	p[0] = (uint8)(x>>2);
	p[1] = (uint8)((x<<6) | ((y>>4)&0x3f));
	p[2] = (uint8)((y<<4) | (dir&0xf));
*/
