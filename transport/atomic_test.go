package transport

import (
	"fmt"
	"math/rand"
	"os"
	"periph.io/x/conn/v3/physic"
	"testing"
	"time"
)

func BenchmarkAtomicFloat64Transport_Access(b *testing.B) {
	var intensityTransport AtomicFloat64Transport
	done := make(chan bool)
	/*go func() {
		// publisher
		ticker := time.NewTicker((10 * physic.KiloHertz).Period())
		for {
			select {
			case <-ticker.C:
				intensityTransport.Update(rand.Float64())
				//fmt.Println("updated")
			case <-done:
				return
			}
		}
	}()*/
	go publishRandom(&intensityTransport, done)
	if err := intensityTransport.EnsureReady(time.Second, time.Millisecond); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	for n := 0; n < b.N; n++ {
		_ = intensityTransport.Access()
		//fmt.Println("accessed")
	}
	done <- true
}

func TestAtomicFloat64Transport_Access(t *testing.T) {
	var intensityTransport AtomicFloat64Transport
	done := make(chan bool)
	go func() {
		// publisher
		ticker := time.NewTicker((10 * physic.KiloHertz).Period())
		for {
			select {
			case <-ticker.C:
				intensityTransport.Update(rand.Float64())
			case <-done:
				return
			}
		}
	}()
	if err := intensityTransport.EnsureReady(time.Second, time.Millisecond); err != nil {
		t.Error(err)
	}
	for n := 0; n < 100; n++ {
		_ = intensityTransport.Access()
	}
	done <- true
}
