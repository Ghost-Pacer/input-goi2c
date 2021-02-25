package main

import (
	"flag"
	"github.com/Ghost-Pacer/input-goi2c/pkg/bno055"
	"gopkg.in/tomb.v2"
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
var Verbose = flag.Bool("v", false, "print all sent messages (do not use in production)")

func printv(vals ...interface{}) {
	if *Verbose {
		log.Println(vals...)
	}
}

func run() error {
	mainTomb := new(tomb.Tomb)

	if _, err := host.Init(); err != nil {
		return err
	}
	log.Println("periph: initted host")

	bus, err := i2creg.Open(*I2CBus)
	if err != nil {
		return err
	}
	log.Println("periph: initted bus")
	defer bus.Close()

	bno, err := bno055.New(bus, uint16(*I2CAddr))
	log.Println("bno055: initted bno055")
	if err != nil {
		return err
	}
	defer bno.Halt()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	mainTomb.Go(func() error {
		return doWork(mainTomb, bno)
	})

	caughtSignal := <-signalChannel
	log.Println("Caught", caughtSignal, "shutting down...")
	mainTomb.Kill(nil)
	return mainTomb.Wait()
}

func doWork(tomb *tomb.Tomb, bno *bno055.Dev) error {
	ticker := time.NewTicker(*RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tomb.Dying():
			return nil
		case <-ticker.C:
			printv("ticked")

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
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
