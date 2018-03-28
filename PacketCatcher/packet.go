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
			_, err := bufreader.Read(buf)
			if err != nil {
				break
			}
			r := bytes.NewBuffer(buf)
			for {
				if r.Len() == 0 {
					break
				}
				//TODO: Stop using string and transfer to hex
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
				} else if ok && item < 0 {
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
				if takeLength < 0 {
					//fmt.Printf("Length of Packet is less than 0. Breaking: %v", takeLength)
					break
				}
				body := r.Next(takeLength)
				packetc <- New(packetID, body, isPacketSent)
			}
		}
	}()
	return packetc
}

//ParseChannelV2 should be an updated version with refactored logic
func ParseChannelV2(conn io.Reader) chan *Packet {
	var isSent bool
	packetc := make(chan *Packet, 0)
	go func() {
		r := bufio.NewReader(conn)
		rBuffer1 := make([]byte, 2048)
		errCount := 0
		for {
			n, err := r.Read(rBuffer1)
			if err != nil {
				errCount++
				if errCount > 5 {
					//fmt.Println("Have more than 5 errors while reading from connection. Exiting.")
					close(packetc)
					break
				}
				//fmt.Printf("Got an error while reading from connection: %v\n", err)
				continue
			}
			rBuffer := rBuffer1[:n]
			next := true
			//fmt.Printf("Incoming packet is:[%d] %x\n", n, rBuffer)
			if n < 2 {
				//fmt.Println("Array length is less then 2. Wrong packet. Discarding")
				continue
			}
			dataReader := bytes.NewBuffer(rBuffer)
			for dataReader.Len() > 0 && next {
				var pid uint16
				if err := binary.Read(dataReader, binary.LittleEndian, &pid); err != nil {
					errCount++
					//fmt.Printf("Got an error while decoding dataReader: %v\n", err)
					next = false
					continue
				}
				if pid == 0 {
					//fmt.Println("Found strange packet which ID is 0x0000. Skiping")
					continue
				}
				switch pid {
				case 0x5252:
					isSent = false
					continue
				case 0x5353:
					isSent = true
					continue
				case 0x592c:
					continue
				case 0x007f:
					continue
				case 0x09a1:
					continue
				default:
				}
				pidString := fmt.Sprintf("%04x", pid)
				//fmt.Printf("Found packetID: 0x%s\n", pidString)
				if dataReader.Len() == 0 {
					//fmt.Println("Packet doesn't contain any data.")
					packetc <- New(pidString, make([]byte, 0), isSent)
					continue
				}
				var plen int16
				if l, ok := packetLength[pidString]; !ok || l == -1 {
					//fmt.Printf("Present in map: %v\tLength: %d\nVariable Length. Extracting...", ok, l)
					if err := binary.Read(dataReader, binary.LittleEndian, &plen); err != nil {
						errCount++
						//fmt.Printf("Got an error while decoding dataReader: %v\n", err)
						next = false
						continue
					}
					//fmt.Printf("Got packet length: %d\n", plen)
					plen -= 2
				} else if ok {
					//fmt.Printf("Present in map: [%v] Length: [%d] Good!\n", ok, l)
					plen = int16(l)
				} else if l == 0 {
					//fmt.Println("Found zero-body packet")
					packetc <- New(pidString, make([]byte, 0), isSent)
					continue
				}
				plen -= 2
				//fmt.Printf("Rest packetlength is: [%d]. dataReader length is: [%d]\n", plen, dataReader.Len())
				if plen < 0 {
					packetc <- New(pidString, make([]byte, 0), isSent)
					continue
				}
				if plen > 0 && plen <= int16(dataReader.Len()) {
					packetData := dataReader.Next(int(plen)) //TODO: Look carefuly at this cauz maybe reader can't perform zero-len reading
					//fmt.Printf("Prepared packetBody: %x\n", packetData)
					packetc <- New(pidString, packetData, isSent)

				}
			}

		}
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
