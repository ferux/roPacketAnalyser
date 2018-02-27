package PacketTypes

import (
	"bytes"
	"log"

	"github.com/ferux/roPacketAnalyser/PacketCatcher"
)

//Replace assignments use ^\s*\w\[\"(\w+)\"\]\s\=\sBytesTo(\w+).*$
//To p.$1 = m["$1"].($2)

//Replace ^\s*(\w+)\s(\w+)
//TO: m["$1"] = BytesTo$2(b)

//PacketTypes Maps the functions to packetID
var packetTypes map[string]func(*PacketCatcher.Packet) map[string]interface{}

func init() {
	packetTypes = make(map[string]func(*PacketCatcher.Packet) map[string]interface{}, 0)
	packetTypes["0230"] = packet0230ToMap
	packetTypes["08c8"] = packet08c8ToMap
	packetTypes["02e1"] = packet02e1ToMap
	packetTypes["0064"] = packet0064ToMap
	packetTypes["0081"] = packet0081ToMap //PACKET_SC_NOTIFY_BAN
	packetTypes["07db"] = packet07dbToMap //PACKET_ZC_HO_PAR_CHANGE
	packetTypes["0080"] = packet0080ToMap //PACKET_ZC_NOTIFY_VANISH
	packetTypes["07f6"] = packet07f6ToMap //PACKET_ZC_NOTIFY_EXP
	packetTypes["022e"] = packet022eToMap //PACKET_ZC_PROPERTY_HOMUN
	packetTypes["00b0"] = packet00b0ToMap //PACKET_ZC_PAR_CHANGE
	log.Printf("[packetTypes ] Successfuly imported %6d rows\n", len(packetTypes))
}

//MakeMap for specified packetID
func MakeMap(p *PacketCatcher.Packet) map[string]interface{} {
	m, ok := packetTypes[p.PacketID]
	if !ok {
		return nil
	}
	return m(p)
}

func packet0230ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)

	m["PacketID"] = p.PacketID
	m["type"] = BytesToInt8(b)
	m["state"] = BytesToInt8(b)
	m["GID"] = BytesToInt32(b)
	m["data"] = BytesToInt32(b)
	return m
	/*
		this+0x0 short PacketID;
		this+0x2 char type;
		this+0x3 char state;
		this+0x4 int GID;
		this+0x8 int data;
	*/
}

func packet08c8ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["GID"] = BytesToUint32(b)
	m["targetGID"] = BytesToUint32(b)
	m["startTime"] = BytesToUint32(b)
	m["attackMT"] = BytesToInt32(b)
	m["attackedMT"] = BytesToInt32(b)
	m["damage"] = BytesToInt32(b)
	m["IsSPDamage"] = BytesToUint8(b)
	m["count"] = BytesToInt16(b)
	m["action"] = BytesToUint8(b)
	m["leftDamage"] = BytesToInt32(b)
	return m
}

func packet02e1ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["GID"] = BytesToUint32(b)
	m["targetGID"] = BytesToUint32(b)
	m["startTime"] = BytesToUint32(b)
	m["attackMT"] = BytesToInt32(b)
	m["attackedMT"] = BytesToInt32(b)
	m["damage"] = BytesToInt32(b)
	m["IsSPDamage"] = BytesToUint8(b)
	m["count"] = BytesToInt16(b)
	m["action"] = BytesToUint8(b)
	m["leftDamage"] = BytesToInt32(b)
	return m
}

func packet0064ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["Version"] = BytesToUint32(b)
	m["ID"] = BytesToByteArray(b, 24)
	m["Passwd"] = BytesToByteArray(b, 24)
	m["clienttype"] = BytesToUint8(b)

	/*
		this+0x0 short PacketID;
		this+0x2 unsigned long Version;
		this+0x6 unsigned char ID[24];
		this+0x1e unsigned char Passwd[24];
		this+0x36 unsigned char clienttype;

			PacketID   string
			Version    uint32
			ID         [24]byte
			Passwd     [24]byte
			ClientType uint8
	*/
	return m
}

func packet0081ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["ErrorCode"] = BytesToUint8(b)
	/*
		this+0x0 short PacketID;
		this+0x2 unsigned char ErrorCode;
	*/
	return m
}

func packet07dbToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["var"] = BytesToUint16(b)
	m["value"] = BytesToInt32(b)
	/*
		this+0x0 short PacketID;
		this+0x2 unsigned short var;
		this+0x4 int value;
	*/
	return m
}

func packet0080ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["GID"] = BytesToUint32(b)
	m["type"] = BytesToUint8(b)
	/*
		PacketID string
		GID      uint32
		Type     uint8
	*/
	return m
}

func packet07f6ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["AID"] = BytesToUint32(b)
	m["Amount"] = BytesToInt32(b)
	m["VarID"] = BytesToUint16(b)
	m["ExpType"] = BytesToInt16(b)
	return m
}

func packet022eToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["SzName"] = BytesToByteArray(b, 24)
	m["BModified"] = BytesToUint8(b)
	m["NLevel"] = BytesToInt16(b)
	m["NFullness"] = BytesToInt16(b)
	m["NRelationship"] = BytesToInt16(b)
	m["ITID"] = BytesToUint16(b)
	m["Atk"] = BytesToInt16(b)
	m["Matk"] = BytesToInt16(b)
	m["Hit"] = BytesToInt16(b)
	m["Critical"] = BytesToInt16(b)
	m["Def"] = BytesToInt16(b)
	m["Mdef"] = BytesToInt16(b)
	m["Flee"] = BytesToInt16(b)
	m["Aspd"] = BytesToInt16(b)
	m["HP"] = BytesToInt16(b)
	m["MaxHP"] = BytesToInt16(b)
	m["SP"] = BytesToInt16(b)
	m["MaxSP"] = BytesToInt16(b)
	m["Exp"] = BytesToInt32(b)
	m["MaxEXP"] = BytesToInt32(b)
	m["SKPoint"] = BytesToInt16(b)
	m["ATKRange"] = BytesToInt16(b)
	return m
}

func packet00b0ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["varID"] = BytesToUint16(b)
	m["count"] = BytesToInt32(b)
	return m
}

func packet0069ToMap(p *PacketCatcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	// b := bytes.NewBuffer(p.Body)
	// m["PacketID"] = p.PacketID

	return m
}

/*
//Populate struct from map
func (p *Packet0069) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet00B0 type")
	}
		PacketID      string
	PacketLength  int16
	AuthCode      int32
	AID           uint32
	UserLevel     uint32
	LastLoginIP   uint32
	LastLoginTime [26]byte
	Sex           uint8
	return nil
}
*/
