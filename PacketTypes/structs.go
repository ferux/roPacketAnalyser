package PacketTypes

import "errors"

//PacketType is a interface for populating struct from map
type PacketType interface {
	Populate(map[string]interface{}) error
}

var (
	_ PacketType = (*Packet0230)(nil)
	_ PacketType = (*Packet08c8)(nil)
	_ PacketType = (*Packet02e1)(nil)
	_ PacketType = (*Packet0064)(nil)
	_ PacketType = (*Packet0081)(nil)
	_ PacketType = (*Packet07db)(nil)
	_ PacketType = (*Packet0080)(nil)
	_ PacketType = (*Packet07F6)(nil)
	_ PacketType = (*Packet022E)(nil)
	_ PacketType = (*Packet00B0)(nil)
	_ PacketType = (*Packet0069)(nil)
	_ PacketType = (*Packet09DD)(nil)
)

//Replace assignments use ^\s*\w\[\"(\w+)\"\]\s\=\sBytesTo(\w+).*$
//To p.$1 = m["$1"].($2)

//To convert struct to populate:
//Replace: ^\s*([\w_]+)\s([\w\[\]]+)$
//     To: p.$1 = m["$1"].($2)

//Packet0230 represents 0x0230 packet
type Packet0230 struct {
	PacketID string
	Type     int8
	State    int8
	GID      int32
	Data     int32
}

//Populate struct from map
func (p *Packet0230) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet0230 type. It is " + packetID)
	}
	p.PacketID = packetID
	p.State = m["state"].(int8)
	p.Type = m["type"].(int8)
	p.Data = m["data"].(int32)
	p.GID = m["GID"].(int32)
	return nil
}

//Packet08c8 represents 0x08c8 packet
type Packet08c8 struct {
	PacketID   string
	GID        uint32
	TargetGID  uint32
	StartTime  uint32
	AttackMT   int32
	AttackedMT int32
	Damage     int32
	IsSPDamage uint8
	Count      int16
	Action     uint8
	LeftDamage int32
}

//Populate struct from map
func (p *Packet08c8) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet08c8 type")
	}
	p.PacketID = packetID
	p.GID = m["GID"].(uint32)
	p.TargetGID = m["targetGID"].(uint32)
	p.StartTime = m["startTime"].(uint32)
	p.AttackMT = m["attackMT"].(int32)
	p.AttackedMT = m["attackedMT"].(int32)
	p.Damage = m["damage"].(int32)
	p.IsSPDamage = m["IsSPDamage"].(uint8)
	p.Count = m["count"].(int16)
	p.Action = m["action"].(uint8)
	p.LeftDamage = m["leftDamage"].(int32)
	return nil
}

//Packet02e1 represents packet 0x02e1
type Packet02e1 struct {
	PacketID   string
	GID        uint32
	TargetGID  uint32
	StartTime  uint32
	AttackMT   int32
	AttackedMT int32
	Damage     int32
	IsSPDamage uint8
	Count      int16
	Action     uint8
	LeftDamage int32
}

//Populate struct from map
func (p *Packet02e1) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet02e1 type")
	}
	p.PacketID = packetID
	p.GID = m["GID"].(uint32)
	p.TargetGID = m["targetGID"].(uint32)
	p.StartTime = m["startTime"].(uint32)
	p.AttackMT = m["attackMT"].(int32)
	p.AttackedMT = m["attackedMT"].(int32)
	p.Damage = m["damage"].(int32)
	p.IsSPDamage = m["IsSPDamage"].(uint8)
	p.Count = m["count"].(int16)
	p.Action = m["action"].(uint8)
	p.LeftDamage = m["leftDamage"].(int32)
	return nil
}

//Packet0064 represents packet 0x0064
type Packet0064 struct {
	PacketID   string
	Version    uint32
	ID         [24]byte
	Passwd     [24]byte
	ClientType uint8
}

//Populate struct from map
func (p *Packet0064) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet0064 type")
	}
	p.PacketID = packetID
	p.Version = m["Version"].(uint32)
	p.ID = m["ID"].([24]byte)
	p.Passwd = m["Passwd"].([24]byte)
	p.ClientType = m["clienttype"].(uint8)
	return nil
}

//Packet0081 represents 0x81 packet
type Packet0081 struct {
	PacketID  string
	ErrorCode uint8
}

//Populate struct from map
func (p *Packet0081) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet0081 type")
	}
	p.PacketID = packetID
	p.ErrorCode = m["ErrorCode"].(uint8)
	return nil
}

//Packet07db represents 0x07db packet
type Packet07db struct {
	PacketID string
	Var      uint16
	Value    int32
	/*
		this+0x0 short PacketID;
		this+0x2 unsigned short var;
		this+0x4 int value;
	*/
}

//Populate struct from map
func (p *Packet07db) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet07db type")
	}
	p.PacketID = packetID
	p.Var = m["var"].(uint16)
	p.Value = m["value"].(int32)
	return nil
}

//Packet0080 PACKET_ZC_NOTIFY_VANISH
type Packet0080 struct {
	PacketID string
	GID      uint32
	Type     uint8
}

//Populate struct from map
func (p *Packet0080) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet0080 type")
	}
	p.PacketID = packetID
	p.GID = m["GID"].(uint32)
	p.Type = m["type"].(uint8)
	return nil
}

//Packet07F6 PACKET_ZC_NOTIFY_EXP
type Packet07F6 struct {
	PacketID string
	AID      uint32
	Amount   int32
	VarID    uint16
	ExpType  int16
}

//Populate struct from map
func (p *Packet07F6) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not PacketF607 type")
	}
	p.PacketID = packetID
	p.AID = m["AID"].(uint32)
	p.Amount = m["Amount"].(int32)
	p.VarID = m["VarID"].(uint16)
	p.ExpType = m["ExpType"].(int16)
	return nil
	/*
				m["PacketID"] = p.PacketID
		p.AID = m["AID"].(Uint32)
		p.Amount = m["Amount"].(Int32)
		p.VarID = m["VarID"].(Uint16)
		p.ExpType = m["ExpType"].(Int16)
	*/
}

//Packet022E PACKET_ZC_PROPERTY_HOMUN
type Packet022E struct {
	PacketID      string
	SzName        []byte
	BModified     uint8
	NLevel        int16
	NFullness     int16
	NRelationship int16
	ITID          uint16
	Atk           int16
	Matk          int16
	Hit           int16
	Critical      int16
	Def           int16
	Mdef          int16
	Flee          int16
	Aspd          int16
	HP            int16
	MaxHP         int16
	SP            int16
	MaxSP         int16
	Exp           int32
	MaxEXP        int32
	SKPoint       int16
	ATKRange      int16
}

//Populate struct from map
func (p *Packet022E) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not PacketF607 type")
	}
	p.PacketID = packetID
	p.SzName = m["SzName"].([]byte)
	p.BModified = m["BModified"].(uint8)
	p.NLevel = m["NLevel"].(int16)
	p.NFullness = m["NFullness"].(int16)
	p.NRelationship = m["NRelationship"].(int16)
	p.ITID = m["ITID"].(uint16)
	p.Atk = m["Atk"].(int16)
	p.Matk = m["Matk"].(int16)
	p.Hit = m["Hit"].(int16)
	p.Critical = m["Critical"].(int16)
	p.Def = m["Def"].(int16)
	p.Mdef = m["Mdef"].(int16)
	p.Flee = m["Flee"].(int16)
	p.Aspd = m["Aspd"].(int16)
	p.HP = m["HP"].(int16)
	p.MaxHP = m["MaxHP"].(int16)
	p.SP = m["SP"].(int16)
	p.MaxSP = m["MaxSP"].(int16)
	p.Exp = m["Exp"].(int32)
	p.MaxEXP = m["MaxEXP"].(int32)
	p.SKPoint = m["SKPoint"].(int16)
	p.ATKRange = m["ATKRange"].(int16)
	/*
			   m["SzName"] =  BytesToByteArray(b, 24)
		p.BModified = m["BModified"].(uint8)
		p.NLevel = m["NLevel"].(int16)
		p.NFullness = m["NFullness"].(int16)
		p.NRelationship = m["NRelationship"].(int16)
		p.ITID = m["ITID"].(uint16)
		p.Atk = m["Atk"].(int16)
		p.Matk = m["Matk"].(int16)
		p.Hit = m["Hit"].(int16)
		p.Critical = m["Critical"].(int16)
		p.Def = m["Def"].(int16)
		p.Mdef = m["Mdef"].(int16)
		p.Flee = m["Flee"].(int16)
		p.Aspd = m["Aspd"].(int16)
		p.HP = m["HP"].(int16)
		p.MaxHP = m["MaxHP"].(int16)
		p.SP = m["SP"].(int16)
		p.MaxSP = m["MaxSP"].(int16)
		p.Exp = m["Exp"].(int32)
		p.MaxEXP = m["MaxEXP"].(int32)
		p.SKPoint = m["SKPoint"].(int16)
		p.ATKRange = m["ATKRange"].(int16)
	*/
	return nil
}

// packet 0x22e
//struct PACKET_ZC_PROPERTY_HOMUN

//Packet00B0 PACKET_ZC_PAR_CHANGE
type Packet00B0 struct {
	PacketID string
	VarID    uint16
	Count    int32
}

//Populate the packet
func (p *Packet00B0) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet00B0 type")
	}
	p.PacketID = packetID
	p.VarID = m["varID"].(uint16)
	p.Count = m["count"].(int32)
	return nil
}

/*
// packet 0xb0 -- 0x00b0
struct PACKET_ZC_PAR_CHANGE {
	this+0x0  short PacketID;
	this+0x2  unsigned short varID;
	this+0x4  int count;
};
*/
/*
'0069' => [
	'account_server_info',
	'x2 a4 a4 a4 a4 a26 C a*',
	[qw(sessionID accountID sessionID2 lastLoginIP lastLoginTime accountSex serverInfo)]],

struct PACKET_AC_ACCEPT_LOGIN {
	this+0x0 short PacketID;
	this+0x2 short PacketLength;
	this+0x4 int AuthCode;
	this+0x8 unsigned long AID;
	this+0xc unsigned long userLevel;
	this+0x10 unsigned long lastLoginIP;
	this+0x14 char lastLoginTime[26];
	this+0x2e unsigned char Sex;
	this+0x2f struct SERVER_ADDR ServerList[];
};
*/
//Packet 0069 PACKET_AC_ACCEPT_LOGIN
type Packet0069 struct {
	PacketID      string
	PacketLength  int16
	AuthCode      int32
	AID           uint32
	UserLevel     uint32
	LastLoginIP   uint32
	LastLoginTime [26]byte
	Sex           uint8
	//todo: server_addr if needed
}

//Populate struct from map
func (p *Packet0069) Populate(m map[string]interface{}) error {
	packetID, ok := m["PacketID"].(string)
	if !ok {
		return errors.New("passed map is not Packet00B0 type")
	}
	p.PacketID = packetID
	p.PacketLength = m["PacketLength"].(int16)
	p.AuthCode = m["AuthCode"].(int32)
	p.AID = m["AID"].(uint32)
	p.UserLevel = m["UserLevel"].(uint32)
	p.LastLoginIP = m["LastLoginIP"].(uint32)
	p.LastLoginTime = m["LastLoginTime"].([26]byte)
	p.Sex = m["Sex"].(uint8)
	return nil
}

//Packet09DD -- ACTOR_EXISTS
type Packet09DD struct {
	// PacketLength uint16
	ObjectType   uint8
	ID           uint32
	CharID       uint32
	WalkSpeed    int16
	Opt1         int16
	Opt2         int16
	Option       int32
	Type         int16
	HairStyle    int16
	Weapon       int16
	Shield       int16
	LowHead      int16
	TopHead      int16
	MidHead      int16
	HairColor    int16
	ClothesColor int16
	HeadDir      int16
	Costume      int16
	GuildID      uint32
	EmblemID     uint16
	Manner       int16
	Opt3         int32
	Stance       uint8
	Sex          uint8
	CoordX       int16
	CoordY       int16
	Direction    uint8
	// Coords string
	XSize uint8
	YSize uint8
	Act   uint8
	Lv    int16
	Font  int16
	Opt4  string
	Name  string
}

//Populate the packet
func (p *Packet09DD) Populate(m map[string]interface{}) error {
	p.ObjectType = m["ObjectType"].(uint8)
	p.ID = m["ID"].(uint32)
	p.CharID = m["CharID"].(uint32)
	p.WalkSpeed = m["WalkSpeed"].(int16)
	p.Opt1 = m["Opt1"].(int16)
	p.Opt2 = m["Opt2"].(int16)
	p.Option = m["Option"].(int32)
	p.Type = m["Type"].(int16)
	p.HairStyle = m["HairStyle"].(int16)
	p.Weapon = m["Weapon"].(int16)
	p.Shield = m["Shield"].(int16)
	p.LowHead = m["LowHead"].(int16)
	p.TopHead = m["TopHead"].(int16)
	p.MidHead = m["MidHead"].(int16)
	p.HairColor = m["HairColor"].(int16)
	p.ClothesColor = m["ClothesColor"].(int16)
	p.HeadDir = m["HeadDir"].(int16)
	p.Costume = m["Costume"].(int16)
	p.GuildID = m["GuildID"].(uint32)
	p.EmblemID = m["EmblemID"].(uint16)
	p.Manner = m["Manner"].(int16)
	p.Opt3 = m["Opt3"].(int32)
	p.Stance = m["Stance"].(uint8)
	p.Sex = m["Sex"].(uint8)
	coords := m["Coords"].([]uint8)
	p.CoordX, p.CoordY, p.Direction = BytesToXYDir(coords)
	p.XSize = m["XSize"].(uint8)
	p.YSize = m["YSize"].(uint8)
	p.Act = m["Act"].(uint8)
	p.Lv = m["Lv"].(int16)
	p.Font = m["Font"].(int16)
	p.Opt4 = m["Opt4"].(string)
	p.Name = m["Name"].(string)
	return nil
}
