package bno055

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const (
	ChipIDWord = -1120 // 0xFB.A0 (registers 0x01.0x00) interpreted as signed int16
	ModeAddr   = 0x3D
	ModeConfig = 0x00
	ModeNDOF   = 0x0C
)

const (
	QuatScale        = float32(1.0 / (1 << 14))
	EulerScale       = float32(1.0 / 16)
	LinearAccelScale = float32(1.0 / 100)
)

// addresses of LSBs, MSBs assumed +0x01
const (
	QuatAddr        = 0x20
	EulerAddr       = 0x1A
	LinearAccelAddr = 0x28
)

type Dev struct {
	transport conn.Conn
}

func New(bus i2c.Bus, addr uint16) (*Dev, error) {
	dev := &Dev{
		transport: &i2c.Dev{Bus: bus, Addr: addr},
	}
	// TODO power up, load calibration
	if err := dev.writeByte(ModeAddr, ModeConfig); err != nil {
		return nil, err
	}
	time.Sleep(20 * time.Millisecond)

	if err := dev.writeByte(ModeAddr, ModeNDOF); err != nil {
		return nil, err
	}
	time.Sleep(20 * time.Millisecond)

	chipID, err := dev.readWords(0x00, 1)
	if err != nil {
		return nil, err
	}
	if chipID[0] != ChipIDWord {
		return nil, fmt.Errorf("expected chip ID %x but got %x", ChipIDWord, chipID)
	}

	return dev, nil
}

func (dev *Dev) writeByte(addr byte, value byte) error {
	return dev.transport.Tx([]byte{addr, value}, nil)
}

func (dev *Dev) readWords(startAddr byte, numWords int) ([]int16, error) {
	buf := make([]byte, 2*numWords)
	if err := dev.transport.Tx([]byte{startAddr}, buf); err != nil {
		return nil, err
	}
	out := make([]int16, numWords)
	for i := 0; i < numWords; i++ {
		lsb := i * 2
		msb := lsb + 1
		out[i] = int16(buf[msb])<<8 | int16(buf[lsb])
	}
	return out, nil
}

func (dev *Dev) readScaledValues(startAddr byte, numWords int, scaleFactor float32) ([]float32, error) {
	out := make([]float32, numWords)
	unscaledWords, err := dev.readWords(startAddr, numWords)
	if err != nil {
		return nil, err
	}
	for i := 0; i < numWords; i++ {
		out[i] = float32(unscaledWords[i]) * scaleFactor
	}
	return out, nil
}

func (dev *Dev) ReadQuat() ([4]float32, error) {
	var dst [4]float32
	src, err := dev.readScaledValues(QuatAddr, 4, QuatScale)
	if err != nil {
		return [4]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) ReadEuler() ([3]float32, error) {
	var dst [3]float32
	src, err := dev.readScaledValues(EulerAddr, 3, EulerScale)
	if err != nil {
		return [3]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) ReadLinearAccel() ([3]float32, error) {
	var dst [3]float32
	src, err := dev.readScaledValues(LinearAccelAddr, 3, LinearAccelScale)
	if err != nil {
		return [3]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) Halt() error {
	// TODO power down / maybe save config
	log.Println("Halted bno055.")
	return nil
}
