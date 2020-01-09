package catcher

import (
	"bytes"
	"fmt"
	"strconv"
)

func parseHexInt32(r *bytes.Buffer) int32 {
	data := reverse(r.Next(4))
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0
	}
	return int32(v)
}

func parseHexUInt32(r *bytes.Buffer) uint32 {
	data := reverse(r.Next(4))
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

func parseHexInt16(r *bytes.Buffer) int16 {
	data := reverse(r.Next(2))
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0
	}
	return int16(v)
}

func parseHexUInt16(r *bytes.Buffer) uint16 {
	data := reverse(r.Next(2))
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0
	}
	return uint16(v)
}

func parseHexInt8(r *bytes.Buffer) int8 {
	data := r.Next(1)
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return 0
	}
	return int8(v)
}

func parseHexUInt8(r *bytes.Buffer) uint8 {
	data := r.Next(1)
	s := fmt.Sprintf("0x%x", data)
	v, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0
	}
	return uint8(v)
}

func reverse(data []byte) []byte {
	newData := make([]byte, 0)
	for i := len(data) - 1; i >= 0; i-- {
		newData = append(newData, data[i])
	}
	return newData
}
