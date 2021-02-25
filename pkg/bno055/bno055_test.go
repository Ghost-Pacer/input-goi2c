package bno055

import (
	"periph.io/x/conn/v3/i2c/i2creg"
	"testing"
)

func BenchmarkCycle(b *testing.B) {
	bus, err := i2creg.Open("2")
	if err != nil {
		b.Fatal(err)
	}
	b.Log("periph: initted bus")
	defer bus.Close()

	bno, err := New(bus, 0x28)
	if err != nil {
		b.Fatal(err)
	}
	b.Log("bno055: initted bno055")
	defer bno.Halt()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := bno.ReadQuat()
		if err != nil {
			b.Fatal(err)
		}

		_, err = bno.ReadLinearAccel()
		if err != nil {
			b.Fatal(err)
		}
	}
}
