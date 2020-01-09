package catcher

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	packetLength = make(map[string]int)
	populatePacketLength()
}

const packetLengthFileName = "data/pl.ini"

var packetLength map[string]int

func populatePacketLength() {
	f, err := os.Open(packetLengthFileName)
	if err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)

	for {
		buf, _, err := r.ReadLine()
		if err != nil {
			break
		}
		text := string(buf)
		if len(text) < 2 || text[:2] != "0x" {
			continue
		}
		text = strings.Replace(text, "0x", "", -1)
		text = strings.Replace(text, " ", "", -1)
		text = strings.ToLower(text)
		textArray := strings.Split(text, "=")
		length, err := strconv.Atoi(textArray[1])
		if err != nil {
			continue
		}
		packetLength[textArray[0]] = length
	}

	if item, ok := packetLength["02d8"]; !ok || item != 10 {
		log.Printf("PacketLength loaded with error. Got: %d", item)
	}
	log.Printf("[PacketLength] Successfuly imported %6d rows", len(packetLength))
}
