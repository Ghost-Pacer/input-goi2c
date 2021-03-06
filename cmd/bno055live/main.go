package main

import (
	"flag"
	"github.com/Ghost-Pacer/input-goi2c/pkg/bno055"
	"gonum.org/v1/gonum/spatial/r3"
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
var Threshold = flag.Float64("t", 2, "Euler angle difference threshold")

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

	usr3, err := os.OpenFile("/sys/class/leds/beaglebone:green:usr3/brightness", os.O_RDWR|syscall.O_DIRECT, 0600)
	if err != nil {
		return err
	}
	defer usr3.Close()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	mainTomb.Go(func() error {
		return doWork(mainTomb, bno, usr3)
	})

	caughtSignal := <-signalChannel
	log.Println("Caught", caughtSignal, "shutting down...")
	mainTomb.Kill(nil)
	return mainTomb.Wait()
}

func doWork(tomb *tomb.Tomb, bno *bno055.Dev, led *os.File) error {
	ticker := time.NewTicker(*RefreshInterval)
	defer ticker.Stop()

	var lastEul r3.Vec
	for {
		select {
		case <-tomb.Dying():
			return nil
		case <-ticker.C:
			eulArr, err := bno.ReadEuler()
			if err != nil {
				return err
			}

			eul := r3.Vec{eulArr[0], eulArr[1], eulArr[2]}
			diff := r3.Norm(eul.Sub(lastEul))
			lastEul = eul

			printv("current", eul, "last", lastEul, "diff", diff)

			if diff > *Threshold {
				if _, err := led.WriteString("1"); err != nil {
					return err
				}
				time.Sleep(500 * time.Millisecond)
				if _, err := led.WriteString("0"); err != nil {
					return err
				}
			}
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
