package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ferux/roPacketAnalyser/internal/catcher"
	"github.com/ferux/roPacketAnalyser/internal/filter"
	"github.com/ferux/roPacketAnalyser/internal/packet"
	"github.com/ferux/roPacketAnalyser/internal/rpa"
	"github.com/ferux/roPacketAnalyser/internal/rpacontext"
	"github.com/rs/zerolog"

	"github.com/gen2brain/beeep"
)

const laddr = ":13554"

//TODO: Remove vending from this package
func main() {
	const configpath = "config.json"

	settings, err := rpa.ConfigFromFile(configpath)
	if err != nil {
		panic(err)
	}

	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log := zerolog.New(&cw).With().Timestamp().Caller().Logger()

	ctx := rpacontext.NewGracefulContext()
	ctx = rpacontext.WithLogger(ctx, log)

	err = app(ctx, settings)
	if err != nil {
		log.Error().Err(err).Msg("running app")
	}
}

func app(ctx context.Context, cfg rpa.Settings) (err error) {
	var f filter.Packet
	if cfg.InclusiveFilter {
		f = filter.NewAllowList(cfg.SupportedPackets)
	} else {
		f = filter.NewBlockList(cfg.SupportedPackets)
	}

	switch cfg.Mode {
	case rpa.AppModeDialer:
		return startDial(ctx, cfg, f)
	case rpa.AppModeServer:
		return startListen(ctx, cfg, f)
	default:
		return rpa.GeneralError("unknown mode " + cfg.Mode)
	}
}

func startDial(ctx context.Context, cfg rpa.Settings, f filter.Packet) (err error) {
	log := rpacontext.Logger(ctx)

	var conn net.Conn
	for i := 0; i < 10; i++ {
		log.Debug().Str("addr", cfg.Addr).Msg("connecting")

		conn, err = net.Dial("tcp", cfg.Addr)
		if err != nil {
			log.Warn().Err(err).Msg("unable to dial")
			time.Sleep(time.Second * 5)

			continue
		}

		break
	}

	if conn == nil {
		return rpa.GeneralError("unable to connect to server")
	}

	log.Info().Msg("connected")

	handleConnection(ctx, conn, f)

	return nil
}

func startListen(ctx context.Context, cfg rpa.Settings, f filter.Packet) (err error) {
	log := rpacontext.Logger(ctx)

	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Print(err)
		return
	}
	go func() {
		<-ctx.Done()
		_ = l.Close()
	}()

	log.Info().Str("addr", cfg.Addr).Msg("listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accepting connection: %w", err)
		}

		log.Info().Str("remote_addr", conn.RemoteAddr().String()).Msg("connection accepted")

		handleConnection(ctx, conn, f)
	}
}

func handleConnection(ctx context.Context, conn net.Conn, f filter.Packet) {
	var h packetHandler
	var ok bool
	var err error
	var packet *catcher.Packet
	log := rpacontext.Logger(ctx).With().Str("fn", "handleConnection").Logger()
	pc := catcher.ParseChannel(conn)
	for {
		select {
		case packet = <-pc:
		case <-ctx.Done():
			return
		}

		if !f.Allowed(packet.PacketID) {
			log.Debug().Str("packet_id", packet.PacketID).Msg("not allowed")

			continue
		}

		h, ok = handlers[packet.PacketID]
		if !ok {
			err = handlePacket(packet)
		} else {
			err = h.handlePacket(ctx, packet)
		}

		log.Warn().Err(err).Str("packet_id", packet.PacketID).Msg("handler error")
	}
}

type packetHandler interface {
	handlePacket(ctx context.Context, p *catcher.Packet) (err error)
}

var handlers = map[string]packetHandler{
	"0064": loginHandler{},
}

type loginHandler struct{}

func (loginHandler) handlePacket(ctx context.Context, p *catcher.Packet) (err error) {
	m := packet.MakeMap(p)
	if m == nil {
		return
	}

	rpacontext.Logger(ctx).Info().
		Str("type", "login_handler").
		Str("login", m["ID"].(string)).
		Msg("logged in")

	return nil
}

// TODO: refactor this.
func handlePacket(p *catcher.Packet) (err error) {
	m := packet.MakeMap(p)
	if m == nil {
		log.Printf("%s\n", p.String())
		log.Printf("%x\n", p.Body)
	} else {
		switch p.PacketID {
		case "0064":
			login := fmt.Sprintf("%s", m["ID"].([]byte))
			output := fmt.Sprintf("Logged in as %s\n", login)
			log.Println(output)
		case "02e1":
			p := packet.Packet02e1{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			log.Printf("[%d] did %d damage to %d", p.GID, p.Damage, p.TargetGID)
		case "0230":
			p := packet.Packet0230{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			switch p.State {
			case 1:
				etaMin := time.Duration((911-p.Data)*10) * time.Minute
				etaDate := time.Now().Add(etaMin)
				rate := float32(p.Data) / 911.0 * 100.0
				log.Printf("Homunculus intimacy level changed to %d (Left: %d [%5.2f%%]) (ETA: %.0f / %v)", p.Data, 911-p.Data, rate, etaMin.Minutes(), etaDate.Format("02-01-2006 15:04:05"))
			case 2:
				log.Printf("Homunculus hungry level changed to %d\n", p.Data)
				if p.Data < 20 {
					err = beeep.Notify("Homunculus", fmt.Sprintf("Hungry level is %d", p.Data), "ro.png")
					if err != nil {
						log.Printf("unable to notify: %v", err)
					}
				}
			}
		case "08c8":
			p := packet.Packet08c8{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			log.Printf("Character [%d] did %d damage to %d\n", p.GID, p.Damage, p.TargetGID)
		case "0080":
			p := packet.Packet0080{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			log.Printf("Item [%d] of type [%d] vanished", p.GID, p.Type)
		case "07f6":
			// p := packet.Packet07F6{}
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
			p := packet.Packet022E{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			etaMin := time.Duration((911-p.NRelationship)*10) * time.Minute
			etaDate := time.Now().Add(etaMin)
			rate := float32(p.NRelationship) / 911.0 * 100.0
			log.Printf("Homunculus [%s]: Intimacy: %d (%5.2f%%)\tHunger: %d\tETA: %.0fm (~%.2fh)\tDate: %v", p.SzName, p.NRelationship, rate, p.NFullness, etaMin.Minutes(), etaMin.Hours(), etaDate)
		case "00b0":
			p := packet.Packet00B0{}
			if err = p.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			txt := catcher.GetVarID(p.VarID)
			log.Printf("Parameter %s changed to %d", txt, p.Count)
		case "09dd":
			pd := packet.Packet09DD{}
			if err = pd.Populate(m); err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			if pd.Name == "Tomb" {
				notify := fmt.Sprintf("Found Tomb: %d %d\n", pd.CoordX, pd.CoordY)

				// TODO: move to logger
				fmt.Printf("%v\n", notify)

				// TODO: handle error
				_ = beeep.Notify("TOMB", notify, "ro.png")
			}
		default:
			return rpa.ErrHandlerNotFound
		}

	}

	return nil
}

// TODO: move to internal package.
type exp struct {
	timeStart    time.Time
	totalBaseExp int32
	totalJobExp  int32
	monKilled    int

	mu sync.RWMutex
}

func (e *exp) addBaseExp(i int32) {
	e.mu.Lock()
	if e.monKilled == 0 {
		e.timeStart = time.Now()
	}
	if i >= 0 {
		e.monKilled++
	}
	e.totalBaseExp += i
	e.mu.Unlock()
}

func (e *exp) addJobExp(i int32) {
	e.mu.Lock()
	e.totalJobExp += i
	e.mu.Unlock()
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
