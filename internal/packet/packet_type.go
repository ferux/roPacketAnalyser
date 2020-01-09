package packet

import (
	"bytes"
	"log"

	"github.com/ferux/roPacketAnalyser/internal/catcher"
)

//Replace assignments use ^\s*\w\[\"(\w+)\"\]\s\=\sBytesTo(\w+).*$
//To p.$1 = m["$1"].($2)

//Replace ^\s*(\w+)\s(\w+)
//TO: m["$1"] = BytesTo$2(b)

//packet Maps the functions to packetID
var packetTypes map[string]func(*catcher.Packet) map[string]interface{}

func init() {
	packetTypes = make(map[string]func(*catcher.Packet) map[string]interface{})
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
	packetTypes["09dd"] = packet09ddToMap //PACKET_ZC_ACTOR_EXISTS
	log.Printf("[packetTypes ] Successfuly imported %6d rows\n", len(packetTypes))
}

//MakeMap for specified packetID
func MakeMap(p *catcher.Packet) map[string]interface{} {
	m, ok := packetTypes[p.PacketID]
	if !ok {
		return nil
	}
	return m(p)
}

func packet0230ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 5)
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

func packet08c8ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 11)
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

func packet02e1ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 11)
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

func packet0064ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 5)
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

func packet0081ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 2)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["ErrorCode"] = BytesToUint8(b)
	/*
		this+0x0 short PacketID;
		this+0x2 unsigned char ErrorCode;
	*/
	return m
}

func packet07dbToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 3)
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

func packet0080ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 3)
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

func packet07f6ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 5)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["AID"] = BytesToUint32(b)
	m["Amount"] = BytesToInt32(b)
	m["VarID"] = BytesToUint16(b)
	m["ExpType"] = BytesToInt16(b)
	return m
}

func packet022eToMap(p *catcher.Packet) map[string]interface{} {
	b := bytes.NewBuffer(p.Body)
	return map[string]interface{}{
		"PacketID":      p.PacketID,
		"SzName":        BytesToByteArray(b, 24),
		"BModified":     BytesToUint8(b),
		"NLevel":        BytesToInt16(b),
		"NFullness":     BytesToInt16(b),
		"NRelationship": BytesToInt16(b),
		"ITID":          BytesToUint16(b),
		"Atk":           BytesToInt16(b),
		"Matk":          BytesToInt16(b),
		"Hit":           BytesToInt16(b),
		"Critical":      BytesToInt16(b),
		"Def":           BytesToInt16(b),
		"Mdef":          BytesToInt16(b),
		"Flee":          BytesToInt16(b),
		"Aspd":          BytesToInt16(b),
		"HP":            BytesToInt16(b),
		"MaxHP":         BytesToInt16(b),
		"SP":            BytesToInt16(b),
		"MaxSP":         BytesToInt16(b),
		"Exp":           BytesToInt32(b),
		"MaxEXP":        BytesToInt32(b),
		"SKPoint":       BytesToInt16(b),
		"ATKRange":      BytesToInt16(b),
	}
}

func packet00b0ToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 3)
	b := bytes.NewBuffer(p.Body)
	m["PacketID"] = p.PacketID
	m["varID"] = BytesToUint16(b)
	m["count"] = BytesToInt32(b)
	return m
}

// func packet0069ToMap(p *catcher.Packet) map[string]interface{} {
// b := bytes.NewBuffer(p.Body)
// m["PacketID"] = p.PacketID
// panic("not implemented yet")
// }

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
func packet09ddToMap(p *catcher.Packet) map[string]interface{} {
	m := make(map[string]interface{}, 32)
	b := bytes.NewBuffer(p.Body)
	m["ObjectType"] = BytesToUint8(b)
	// m["ID"] = BytesToByteArray(b, 4)
	// m["CharID"] = BytesToByteArray(b, 4)
	m["ID"] = BytesToUint32(b)
	m["CharID"] = BytesToUint32(b)
	m["WalkSpeed"] = BytesToInt16(b)
	m["Opt1"] = BytesToInt16(b)
	m["Opt2"] = BytesToInt16(b)
	m["Option"] = BytesToInt32(b)
	m["Type"] = BytesToInt16(b)
	m["HairStyle"] = BytesToInt16(b)
	m["Weapon"] = BytesToInt16(b)
	m["Shield"] = BytesToInt16(b)
	m["LowHead"] = BytesToInt16(b)
	m["TopHead"] = BytesToInt16(b)
	m["MidHead"] = BytesToInt16(b)
	m["HairColor"] = BytesToInt16(b)
	m["ClothesColor"] = BytesToInt16(b)
	m["HeadDir"] = BytesToInt16(b)
	m["Costume"] = BytesToInt16(b)
	// m["GuildID"] = BytesToByteArray(b, 4)
	// m["EmblemID"] = BytesToByteArray(b, 2)
	m["GuildID"] = BytesToUint32(b)
	m["EmblemID"] = BytesToUint16(b)
	m["Manner"] = BytesToInt16(b)
	m["Opt3"] = BytesToInt32(b)
	m["Stance"] = BytesToUint8(b)
	m["Sex"] = BytesToUint8(b)
	m["Coords"] = BytesToByteArray(b, 3)
	// m["CoordX"] = BytesToUint8(b)
	// m["CoordY"] = BytesToUint8(b)
	// m["CoordZ"] = BytesToUint8(b)
	m["XSize"] = BytesToUint8(b)
	m["YSize"] = BytesToUint8(b)
	m["Act"] = BytesToUint8(b)
	m["Lv"] = BytesToInt16(b)
	m["Font"] = BytesToInt16(b)
	//'09DD' => ['actor_exists', 'a9 Z*', [qw()]],
	m["Opt4"] = string(BytesToByteArray(b, 9))
	m["Name"] = string(BytesToByteArray(b, b.Len()))
	return m
}

/*actor_exists -- 09DD
let packet09DD = new Parser().endianess('little').int16('len').uint8('object_type').string('ID', {length: "len"}).string('charID', {length: "len"}).int16('walk_speed')
.int16('opt1').int16('opt2').int32('option').int16('type').int16('hair_style').int16('weapon').int16('shield').int16('lowhead')
.int16('tophead').int16('midhead').int16('hair_color').int16('clothes_color').int16('head_dir').int16('costume').string('guildID', {length: "len"})
.string('emblemID', {length: "len"}).int16('manner').int32('opt3').uint8('stance').uint8('sex').string('coords', {length: "len"})
.uint8('xSize').uint8('ySize').uint8('act').int16('lv').int16('font')
.string('opt4', {length: "len"})
.string('name', {length: "len"});
*/
