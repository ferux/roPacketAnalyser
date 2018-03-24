package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ferux/roPacketAnalyser/PacketCatcher"
	"github.com/ferux/roPacketAnalyser/PacketTypes"
	"github.com/gen2brain/beeep"
)

const laddr = ":13554"

var wg sync.WaitGroup

//TODO: Remove vending from this package
func main() {
	startListen()
}

func startDial() {
	defer func() {
		wg.Done()
		log.Println("startDial stopped")
	}()
	log.Println("Dialing to packet sniffer")
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", "127.0.0.1:10101")
		if err != nil {
			log.Printf("Can't dial addr. Reason: %v", err)
			time.Sleep(time.Second * 5)
			log.Printf("Connecting again")
		} else {
			break
		}
	}
	if conn == nil {
		log.Fatal("Can't connect to server")
	}
	log.Println("Successfuly connected")
	serveConnV3(conn)
}

func startListen() {
	defer func() {
		log.Println("startListen stopped")
	}()

	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Waiting for connection from RO Client")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error during connection close: %v\n", conn.Close())
		}
		log.Println("Connection from client accepted")
		serveConnV3(conn)
	}
}

func serveConnV3(conn net.Conn) {
	pc := PacketCatcher.ParseChannel(conn)
	for packet := range pc {
		go servePacket(packet)
		// servePacket(packet)
	}
}

func servePacket(p *PacketCatcher.Packet) {
	m := PacketTypes.MakeMap(p)
	if m == nil {
		// log.Printf("%s\n", p.String())
		// log.Printf("%x\n", p.Body)
	} else {
		switch p.PacketID {
		case "0064":
			login := fmt.Sprintf("%s", m["ID"].([]byte))
			output := fmt.Sprintf("Logged in as %s\n", login)
			log.Println(output)
		case "02e1":
			// p := PacketTypes.Packet02e1{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// log.Printf("[%d] did %d damage to %d", p.GID, p.Damage, p.TargetGID)
		case "0230":
			// p := PacketTypes.Packet0230{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// switch p.State {
			// case 1:
			// 	etaMin := time.Duration((911-p.Data)*10) * time.Minute
			// 	etaDate := time.Now().Add(etaMin)
			// 	rate := float32(p.Data) / 911.0 * 100.0
			// 	log.Printf("Homunculus intimacy level changed to %d (Left: %d [%5.2f%%]) (ETA: %.0f / %v)", p.Data, 911-p.Data, rate, etaMin.Minutes(), etaDate.Format("02-01-2006 15:04:05"))
			// case 2:
			// 	log.Printf("Homunculus hungry level changed to %d\n", p.Data)
			// 	if p.Data < 20 {
			// 		beeep.Notify("Homunculus", fmt.Sprintf("Hungry level is %d", p.Data), "ro.png")
			// 	}
			// }
		case "08c8":
			// p := PacketTypes.Packet08c8{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// log.Printf("Character [%d] did %d damage to %d\n", p.GID, p.Damage, p.TargetGID)
		case "0080":
			// p := PacketTypes.Packet0080{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// log.Printf("Item [%d] of type [%d] vanished", p.GID, p.Type)
		case "07f6":
			// p := PacketTypes.Packet07F6{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// expType := ""
			// switch p.VarID {
			// case 1:
			// 	expType = "Base"
			// 	totalExp.addBaseExp(p.Amount)
			// 	log.Printf("[Kills: %5d] Gained [%5d] %s exp. [TotalBaseExp: %7d] [Per Hour: %9.2f]",
			// 		totalExp.getKills(), p.Amount, expType, totalExp.getBaseExp(), totalExp.getBasePerHour())
			// case 2:
			// 	expType = "Job "
			// 	totalExp.addJobExp(p.Amount)
			// 	log.Printf("[Kills: %5d] Gained [%5d] %s exp. [TotalJobExp:  %7d] [Per Hour: %9.2f]",
			// 		totalExp.getKills(), p.Amount, expType, totalExp.getJobExp(), totalExp.getJobPerHour())
			// default:
			// 	expType = string(p.ExpType)
			// }
			// log.Printf("[KillAmount: %d] Gained [%5d] %s exp. [TotalExp: %7d|%7d] [Per Hour: %9.2f|%9.2f]", totalExp.getKills(), p.Amount, expType, totalExp.getBaseExp(), totalExp.getJobExp(), totalExp.getBasePerHour(), totalExp.getJobPerHour())
		case "022e":
			// p := PacketTypes.Packet022E{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// etaMin := time.Duration((911-p.NRelationship)*10) * time.Minute
			// etaDate := time.Now().Add(etaMin)
			// rate := float32(p.NRelationship) / 911.0 * 100.0
			// log.Printf("Homunculus [%s]: Intimacy: %d (%5.2f%%)\tHunger: %d\tETA: %.0fm (~%.2fh)\tDate: %v", p.SzName, p.NRelationship, rate, p.NFullness, etaMin.Minutes(), etaMin.Hours(), etaDate)
		case "00b0":
			// p := PacketTypes.Packet00B0{}
			// if err := p.Populate(m); err != nil {
			// 	log.Printf("Error: %v\n", err)
			// 	return
			// }
			// txt := PacketCatcher.GetVarID(p.VarID)
			// log.Printf("Parameter %s changed to %d", txt, p.Count)
		case "09dd":
			pd := PacketTypes.Packet09DD{}
			if err := pd.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			if pd.Name == "Tomb" {
				notify := fmt.Sprintf("Found Tomb: %d %d\n", pd.CoordX, pd.CoordY)
				fmt.Printf(notify)
				beeep.Notify("TOMB", notify, "ro.png")
			}
		default:
			// log.Printf("Packet: %v\n", m)
		}
	}
}

type exp struct {
	timeStart    time.Time
	totalBaseExp int32
	totalJobExp  int32
	monKilled    int
	sync.RWMutex
}

func (e *exp) addBaseExp(i int32) {
	e.Lock()
	if e.monKilled == 0 {
		e.timeStart = time.Now()
	}
	if i >= 0 {
		e.monKilled++
	}
	e.totalBaseExp += i
	e.Unlock()
}

func (e *exp) addJobExp(i int32) {
	e.Lock()
	e.totalJobExp += i
	e.Unlock()
}

func (e *exp) getBaseExp() int32 {
	return e.totalBaseExp
}

func (e *exp) getJobExp() int32 {
	return e.totalJobExp
}

func (e *exp) getKills() int {
	return e.monKilled
}

func (e *exp) getAvg() []float64 {
	report := []float64{0.0, 0.0}
	report[0] = float64(e.totalBaseExp*1.0) / float64(e.monKilled*1.0)
	report[1] = float64(e.totalJobExp*1.0) / float64(e.monKilled*1.0)
	return report
}

func (e *exp) getAvgBase() float64 {
	return float64(e.totalBaseExp*1.0) / float64(e.monKilled*1.0)
}

func (e *exp) getAvgJob() float64 {
	return float64(e.totalJobExp*1.0) / float64(e.monKilled*1.0)
}

func (e *exp) getPerHour() []float64 {
	report := []float64{0.0, 0.0}
	d := time.Since(e.timeStart)
	if e.totalBaseExp == 0 {
		report[0] = 0
	} else {
		report[0] = float64(e.totalBaseExp) * 60.0 / d.Minutes()
	}
	if e.totalJobExp == 0 {
		report[1] = 0
	} else {
		report[1] = float64(e.totalJobExp) * 60.0 / d.Minutes()
	}

	return report
}

func (e *exp) getBasePerHour() float64 {
	d := time.Since(e.timeStart)
	if e.totalBaseExp == 0 {
		return 0
	}
	return float64(e.totalBaseExp) * 60.0 / d.Minutes()
}

func (e *exp) getJobPerHour() float64 {
	d := time.Since(e.timeStart)
	if e.totalJobExp == 0 {
		return 0
	}
	return float64(e.totalJobExp) * 60.0 / d.Minutes()
}

var totalExp = exp{}
