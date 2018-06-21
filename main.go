package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ferux/roPacketAnalyser/PacketCatcher"
	"github.com/ferux/roPacketAnalyser/PacketTypes"
	"github.com/gen2brain/beeep"
	"github.com/namsral/flag"
)

var logger = log.New()

type configType struct {
	dial      bool
	addr      string
	attempts  int
	debugMode bool
	debugFile string
}

var conf configType

func main() {
	var dial, debugMode bool
	var addr, debugFile string
	var attempts int
	flag.String(flag.DefaultConfigFlagname, "", "Path to config file")

	flag.BoolVar(&dial, "dial", false, "Set application into dialing mode")
	flag.BoolVar(&debugMode, "debug", false, "Turns debug mode on")

	flag.StringVar(&addr, "addr", ":13554", "Listening/dialing address for incoming/outgoing connection")
	flag.StringVar(&debugFile, "output", "", "(Optional) If set the log will be written into the file")

	flag.IntVar(&attempts, "attempts", 10, "Amount of attempts the application will perform before exiting")
	flag.Parse()

	conf.dial = dial
	conf.addr = addr
	conf.attempts = attempts
	conf.debugMode = debugMode
	conf.debugFile = debugFile

	logger.Infof("Config: %+v", conf)
	if conf.debugMode {
		if len(conf.debugFile) != 0 {
		}
	}
	exitc := make(chan os.Signal, 1)
	finish := make(chan struct{})
	signal.Notify(exitc, syscall.SIGTERM, syscall.SIGQUIT)
	if !conf.dial {
		go startListen(conf.addr, finish)
	} else {
		go startDial(conf.addr, finish)
	}
	<-exitc
	finish <- struct{}{}
}

func startDial(daddr string, exitc <-chan struct{}) {
	logger.WithField("Address:", daddr).Infoln("Dialing...")
	var conn net.Conn
	var err error
	exitc2 := make(chan struct{}, 1)
	next := true
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", daddr)
		if err != nil {
			logger.WithError(err).Warning("Can't dial to specified addres. Trying again in 5 seconds...")
			time.Sleep(time.Second * 5)
			logger.Warningln(" Connecting...")
		} else {
			break
		}
	}
	if conn == nil {
		logger.Fatal("Can't connect to the server")
	}
	logger.Infoln("Successfully connected")
	go func() {
		<-exitc
		next = false
		exitc2 <- struct{}{}
	}()
	for next {
		serveConnV3(conn, exitc2)
	}

}

func startListen(laddr string, exitc <-chan struct{}) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		logger.Print(err)
		return
	}
	logger.Infoln("Waiting for connection from RO Client")
	exitc2 := make(chan struct{}, 1)
	next := true
	go func() {
		<-exitc
		next = false
		exitc2 <- struct{}{}
	}()
	for next {
		conn, err := l.Accept()
		if err != nil {
			logger.WithError(err).Errorln("There was error while accepting connection")
			logger.WithError(conn.Close()).Infoln("Closing connection")
		}
		logger.Infoln("Connection from client accepted")

		serveConnV3(conn, exitc2)
		logger.Warningln("Connection has been closed")
	}
}

func serveConnV3(conn net.Conn, exitc <-chan struct{}) {
	pc := PacketCatcher.ParseChannelV2(conn)
serve:
	for {
		select {
		case packet := <-pc:
			if !packet.IsSent {
				go servePacket(packet)
			}
		case <-exitc:
			logger.Infoln("Got exit signal. Exiting...")
			break serve
		}
	}
	logger.Println("ServeConnV3 ended")
}

func servePacket(p *PacketCatcher.Packet) {
	m := PacketTypes.MakeMap(p)
	if m == nil {
		logger.Printf("%s", p.String())
		logger.Printf("%x", p.Body)
	} else {
		switch p.PacketID {
		case "0064":
			login := fmt.Sprintf("%s", m["ID"].([]byte))
			output := fmt.Sprintf("Logged in as %s", login)
			logger.Println(output)
		case "02e1":
			p := PacketTypes.Packet02e1{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			logger.Printf("[%d] did %d damage to %d", p.GID, p.Damage, p.TargetGID)
		case "0230":
			p := PacketTypes.Packet0230{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			switch p.State {
			case 1:
				etaMin := time.Duration((911-p.Data)*10) * time.Minute
				etaDate := time.Now().Add(etaMin)
				rate := float32(p.Data) / 911.0 * 100.0
				logger.Printf("Homunculus intimacy level changed to %d (Left: %d [%5.2f%%]) (ETA: %.0f / %v)", p.Data, 911-p.Data, rate, etaMin.Minutes(), etaDate.Format("02-01-2006 15:04:05"))
			case 2:
				logger.Printf("Homunculus hungry level changed to %d", p.Data)
				if p.Data < 20 {
					beeep.Notify("Homunculus", fmt.Sprintf("Hungry level is %d", p.Data), "ro.png")
				}
			}
		case "08c8":
			p := PacketTypes.Packet08c8{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			logger.Printf("Character [%d] did %d damage to %d", p.GID, p.Damage, p.TargetGID)
		case "0080":
			p := PacketTypes.Packet0080{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			logger.Printf("Item [%d] of type [%d] vanished", p.GID, p.Type)
		case "07f6":
			p := PacketTypes.Packet07F6{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			expType := ""
			switch p.VarID {
			case 1:
				expType = "Base"
				totalExp.addBaseExp(p.Amount)
				logger.Printf("[Kills: %5d] Gained [%5d] %s exp. [TotalBaseExp: %7d] [Per Hour: %9.2f]",
					totalExp.getKills(), p.Amount, expType, totalExp.getBaseExp(), totalExp.getBasePerHour())
			case 2:
				expType = "Job "
				totalExp.addJobExp(p.Amount)
				logger.Printf("[Kills: %5d] Gained [%5d] %s exp. [TotalJobExp:  %7d] [Per Hour: %9.2f]",
					totalExp.getKills(), p.Amount, expType, totalExp.getJobExp(), totalExp.getJobPerHour())
			default:
				expType = string(p.ExpType)
			}
			logger.Printf("[KillAmount: %d] Gained [%5d] %s exp. [TotalExp: %7d|%7d] [Per Hour: %9.2f|%9.2f]", totalExp.getKills(), p.Amount, expType, totalExp.getBaseExp(), totalExp.getJobExp(), totalExp.getBasePerHour(), totalExp.getJobPerHour())
		case "022e":
			p := PacketTypes.Packet022E{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			etaMin := time.Duration((911-p.NRelationship)*10) * time.Minute
			etaDate := time.Now().Add(etaMin)
			rate := float32(p.NRelationship) / 911.0 * 100.0
			logger.Printf("Homunculus [%s]: Intimacy: %d (%5.2f%%)\tHunger: %d\tETA: %.0fm (~%.2fh)\tDate: %v", p.SzName, p.NRelationship, rate, p.NFullness, etaMin.Minutes(), etaMin.Hours(), etaDate)
		case "00b0":
			p := PacketTypes.Packet00B0{}
			if err := p.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			txt := PacketCatcher.GetVarID(p.VarID)
			logger.Printf("Parameter %s changed to %d", txt, p.Count)
		case "09dd":
			pd := PacketTypes.Packet09DD{}
			if err := pd.Populate(m); err != nil {
				logger.Printf("Error: %v", err)
				return
			}
			if pd.Name == "Tomb" {
				notify := fmt.Sprintf("Found Tomb: %d %d", pd.CoordX, pd.CoordY)
				fmt.Printf(notify)
				beeep.Notify("TOMB", notify, "ro.png")
			}
		case "09db":
			pd := PacketTypes.Packet09DB{}
			if err := pd.Populate(m); err != nil {
				logger.WithError(err).WithField("PacketID", p.PacketID).Warningln("Can't parse packet")
				return
			}
			logger.Infof("%#v", pd)
		case "09dc":
			pd := PacketTypes.Packet09DB{}
			if err := pd.Populate(m); err != nil {
				logger.WithError(err).WithField("PacketID", p.PacketID).Warningln("Can't parse packet")
				return
			}
			if pd.Name == "Tomb" {
				notify := fmt.Sprintf("Tomb Disappeared: %d %d", pd.CoordX, pd.CoordY)
				fmt.Printf(notify)
				beeep.Notify("TOMB", notify, "ro.png")
			}
			logger.Infof("%#v", pd)
		default:
			logger.Printf("Packet: %v", m)
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
