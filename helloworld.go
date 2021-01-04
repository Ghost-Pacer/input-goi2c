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

	/*dev := i2c.Dev{Bus: bus, Addr: I2CAddr}
	log.Print("initted")

	// _, err = bmxx80.NewI2C(bus, 0x78, &bmxx80.DefaultOpts)

	buf_out := make([]byte, 1)
	if err := dev.Tx([]byte{byte(0x00)}, buf_out); err != nil {
		return err
	}

	log.Printf("0x%02X\n", buf_out[0])*/

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
				quat, err := bno.ReadQuat()
				if err != nil {
					panic(err)
				}
				log.Print(quat)
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
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "input-goi2c: %s.\n", err)
		os.Exit(1)
	}
}
