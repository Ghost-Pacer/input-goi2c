package bno055

import (
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const (
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
var (
	QuatAddrs        = [4]byte{0x20, 0x22, 0x24, 0x26}
	EulerAddrs       = [3]byte{0x1A, 0x1C, 0x1E}
	LinearAccelAddrs = [3]byte{0x28, 0x2A, 0x2C}
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

	return dev, nil
}

func (dev *Dev) writeByte(addr byte, value byte) error {
	return dev.transport.Tx([]byte{addr, value}, nil)
}

func (dev *Dev) readWord(addr byte) (int16, error) {
	buf := make([]byte, 2)
	if err := dev.transport.Tx([]byte{addr}, buf); err != nil {
		return 0, err
	}
	// fmt.Printf("%x%x ", buf[0], buf[1])
	return int16(buf[1])<<8 | int16(buf[0]), nil
}

func (dev *Dev) readScaledWords(addrs []byte, scaleFactor float32) ([]float32, error) {
	out := make([]float32, len(addrs))
	for i, addr := range addrs {
		unscaledWord, err := dev.readWord(addr)
		if err != nil {
			return nil, err
		}
		scaledValue := float32(unscaledWord) * scaleFactor
		out[i] = scaledValue
	}
	return out, nil
}

func (dev *Dev) ReadQuat() ([4]float32, error) {
	var dst [4]float32
	src, err := dev.readScaledWords(QuatAddrs[:], QuatScale)
	if err != nil {
		return [4]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) ReadEuler() ([3]float32, error) {
	var dst [3]float32
	src, err := dev.readScaledWords(EulerAddrs[:], EulerScale)
	if err != nil {
		return [3]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) ReadLinearAccel() ([3]float32, error) {
	var dst [3]float32
	src, err := dev.readScaledWords(LinearAccelAddrs[:], LinearAccelScale)
	if err != nil {
		return [3]float32{}, err
	}
	copy(dst[:], src)
	return dst, nil
}

func (dev *Dev) Halt() error {
	// TODO power down / maybe save config
	return nil
}
