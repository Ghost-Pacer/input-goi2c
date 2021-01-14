package main

import (
	"flag"
	"fmt"
	"github.com/Ghost-Pacer/input-goi2c/bno055"
	"github.com/zeromq/goczmq"
	"log"
	"os"
	"os/signal"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"syscall"
	"time"
)

var I2CBus = flag.String("b", "2", "I2C bus")
var I2CAddr = flag.Int("a", 0x28, "I2C address")
var RefreshInterval = flag.Duration("r", 8*time.Millisecond, "refresh interval")
var SocketEndpoint = flag.String("e", "tcp://localhost:51101", "ZMQ bound endpoint")
var Verbose = flag.Bool("v", false, "print all sent messages (do not use in production)")

func printv(vals ...interface{}) {
	if *Verbose {
		log.Println(vals...)
	}
}

func mainImpl() error {
	socket, err := goczmq.NewPub(*SocketEndpoint)
	defer socket.Destroy()
	//socket := zmq4.NewPub(context.Background())
	//defer socket.Close()
	// socket.SetOption("CONFLATE", true)

	if err != nil {
		return err
	}
	log.Println("goczmq: listening on", *SocketEndpoint)

	if _, err := host.Init(); err != nil {
		return err
	}
	log.Println("Periph: initted host")

	bus, err := i2creg.Open(*I2CBus)
	if err != nil {
		return err
	}
	log.Println("Periph: initted bus")
	defer bus.Close()

	bno, err := bno055.New(bus, uint16(*I2CAddr))
	log.Println("GP: initted bno055")
	if err != nil {
		return err
	}
	defer bno.Halt()

	ticker := time.NewTicker(*RefreshInterval)
	defer ticker.Stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	var caughtSignal os.Signal

Main:
	for {
		select {
		case caughtSignal = <-signals:
			break Main
		case <-ticker.C:
			printv("ticked")

			start := time.Now()
			if err := socket.SendFrame([]byte("hello world"), goczmq.FlagNone); err != nil {
				panic(err)
			}
			fmt.Println("Raw time on zmq was %v", time.Since(start))

			printv("\tsent on socket")

			quat, err := bno.ReadQuat()
			if err != nil {
				return err
			}
			printv("\tgot quat", quat)

			eul, err := bno.ReadEuler()
			if err != nil {
				return err
			}
			printv("\tgot eul", eul)

			lin, err := bno.ReadLinearAccel()
			if err != nil {
				return err
			}
			printv("\tgot lin", lin)

		}
	}

	log.Println("Caught", caughtSignal, "shutting down...")
	return nil
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}
