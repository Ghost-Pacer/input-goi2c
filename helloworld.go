package main

import (
	"fmt"
	"github.com/Ghost-Pacer/input-goi2c/bno055"
	"log"
	"os"
	"os/signal"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"syscall"
	"time"
)

const (
	I2CBus          = "2"
	I2CAddr         = 0x28
	RefreshInterval = 1000 * time.Millisecond
)

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	bus, err := i2creg.Open(I2CBus)
	if err != nil {
		return err
	}
	log.Println("Initted bus")
	defer bus.Close()

	bno, err := bno055.New(bus, I2CAddr)
	log.Println("Initted bno055")
	if err != nil {
		return err
	}
	defer bno.Halt()

	ticker := time.NewTicker(RefreshInterval)
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
			log.Print("ticked")

			quat, err := bno.ReadQuat()
			if err != nil {
				return err
			}
			log.Print("\tgot quat", quat)

			eul, err := bno.ReadEuler()
			if err != nil {
				return err
			}
			log.Print("\tgot eul", eul)

			lin, err := bno.ReadLinearAccel()
			if err != nil {
				return err
			}
			log.Print("\tgot lin", lin)

		}
	}

	log.Println("Caught", caughtSignal, "shutting down...")
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "input-goi2c: %s.\n", err)
		os.Exit(1)
	}
}
