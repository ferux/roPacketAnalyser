package PacketCatcher

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var isPacketSent bool

//Packet stores all information about packet
type Packet struct {
	PacketID string
	Body     []byte
	IsSent   bool
}

//New creates a new packet with specified parameters
func New(packetID string, body []byte, isSent bool) *Packet {
	return &Packet{
		PacketID: packetID,
		Body:     body,
		IsSent:   isSent,
	}
}

//ParseChannel listens for net.Conn and parses all incoming packets.
func ParseChannel(conn io.Reader) chan *Packet {
	var isPacketSent bool
	packetc := make(chan *Packet, 0)
	buf := make([]byte, 65536)
	bufreader := bufio.NewReader(conn)
	go func() {
		for {
			n, err := bufreader.Read(buf)
			if err != nil {
				break
			}
			r := bytes.NewBuffer(buf[:n])
			for {
				if r.Len() == 0 {
					break
				}
				var packetIDint uint16
				binary.Read(r, binary.LittleEndian, &packetIDint)
				packetID := fmt.Sprintf("%04x", packetIDint)
				switch packetID {
				case "5252":
					isPacketSent = false
					continue
				case "5353":
					isPacketSent = true
					continue
				case "592c":
					continue
				case "007f":
					continue
				}
				takeLength := r.Len()
				item, ok := packetLength[packetID]
				if ok && item > 2 {
					takeLength = item - 2
				} else if ok && item <= 0 {
					var ln uint16
					binary.Read(r, binary.LittleEndian, &ln)
					if ln == 0 {
						continue
					}
					takeLength = int(ln) - 4
				}
				if takeLength >= r.Len() {
					takeLength = r.Len()
				}
				if takeLength == 0 {
					continue
				}
				body := r.Next(takeLength)
				packetc <- New(packetID, body, isPacketSent)
			}
		}
	}()
	return packetc
}

func ParseChannelV2(conn io.Reader) chan *Packet {
	var isPacketSent bool
	packetc := make(chan *Packet, 0)
	// buf := make([]byte, 65536)
	// bufreader := bufio.NewReader(conn)
	go func() {
		p := New("meh", []byte("y"), isPacketSent)
		packetc <- p
	}()
	return packetc
}

func (p *Packet) String() string {
	var name, direction string
	name, ok := packetName[p.PacketID]
	if !ok {
		name = "UNKNOWN_PACKET"
	}
	if p.IsSent {
		direction = "Sent"
	} else {
		direction = "Recv"
	}
	return fmt.Sprintf("%s -> [%s] (PacketID: %s) Length: %d", direction, name, p.PacketID, len(p.Body))
}

//GetVarID from packetVars
func GetVarID(i uint16) string {
	res, ok := packetVars[i]
	if !ok {
		return ""
	}
	return res
}

//Length returns length of the packet
func (p *Packet) Length() int {
	return len(p.Body) + 2
}
