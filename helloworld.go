package main

import (
	"fmt"
	"os"
	"log"
	"periph.io/x/host/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

const I2CBus = "2"
const I2CAddr = 0x28

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	bus, err := i2creg.Open(I2CBus)
	if err != nil {
		return err
	}
	defer bus.Close()

	if p, ok := bus.(i2c.Pins); ok {
		log.Printf("Using pins SCL: %s  SDA: %s", p.SCL(), p.SDA())
	}
	
	dev := i2c.Dev{Bus: bus, Addr: I2CAddr}
	log.Print("initted")

	buf_out := make([]byte, 1)
	if err := dev.Tx([]byte{byte(0x00)}, buf_out); err != nil {
		return err
	}

	log.Printf("0x%02X\n", buf_out[0])
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "input-goi2c: %s.\n", err)
		os.Exit(1)
	}
}
