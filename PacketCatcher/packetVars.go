package PacketCatcher

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"log"
)

func init() {
	packetVars = make(map[uint16]string, 0)
	populatePacketVars()
}

const packetVarsFileName = "data/db.ini"

var packetVars map[uint16]string

func populatePacketVars() {
	f, err := os.Open(packetVarsFileName)
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
		if text != "[Vars]" {
			continue
		}
		break
	}

	for {
		buf, _, err := r.ReadLine()
		if err != nil {
			break
		}
		text := string(buf)
		if len(text) < 2 {
			continue
		}
		text = strings.Replace(text, " ", "", -1)
		textArray := strings.Split(text, "=")
		varID, err := strconv.Atoi(textArray[0])
		if err != nil {
			continue
		}
		packetVars[uint16(varID)] = textArray[1]
	}
	if item, ok := packetVars[10]; !ok || item != "VAR_HAIRCOLOR" {
		log.Printf("PacketVars loaded with error. Got: %s", item)
	}
	log.Printf("[PacketVars  ] Successfuly imported %6d rows", len(packetVars))
}

