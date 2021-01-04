package main

import (
	"fmt"
	"github.com/Ghost-Pacer/input-goi2c/bno055"
	"log"
	"os"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
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

	ticker := time.NewTicker(RefreshInterval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				log.Print("ticked")

				quat, err := bno.ReadQuat()
				if err != nil {
					panic(err)
				}
				log.Print("\tgot quat", quat)

				eul, err := bno.ReadEuler()
				if err != nil {
					panic(err)
				}
				log.Print("\tgot eul", eul)

				lin, err := bno.ReadLinearAccel()
				if err != nil {
					panic(err)
				}
				log.Print("\tgot lin", lin)

			}
		}
	}()

	time.Sleep(10 * time.Second)
	ticker.Stop()
	done <- true
	log.Println("Done")

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "input-goi2c: %s.\n", err)
		os.Exit(1)
	}
}
