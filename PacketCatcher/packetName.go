package PacketCatcher

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func init() {
	packetName = make(map[string]string, 0)
	populatePacketName()
}

const packetNameFileName = "data/pn.ini"

var packetName map[string]string

func populatePacketName() {
	f, err := os.Open(packetNameFileName)
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
		// text = strings.ToLower(text)
		textArray := strings.Split(text, "=")
		// packetName[textArray[0][:4]] = textArray[1]
		packetName[strings.ToLower(textArray[0])] = textArray[1]
	}
	if item, ok := packetName["0970"]; !ok || item != "PACKET_CH_MAKE_CHAR_NOT_STATS" {
		log.Printf("PacketName loaded with error. Got: %s", item)
	}
	log.Printf("[PacketName  ] Successfuly imported %6d rows", len(packetName))
}
