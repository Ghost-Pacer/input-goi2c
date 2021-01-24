package transport

import (
	"fmt"
	"math/rand"
	"os"
	"periph.io/x/conn/v3/physic"
	"testing"
	"time"
)

func BenchmarkAtomicFloat64_PubSub_Access(b *testing.B) {
	trans := NewAtomicFloat64()
	done := make(chan bool)
	pub := trans.PubView()
	sub := trans.SubView()

	go publishRandomPubSub(pub, done)
	if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	for n := 0; n < b.N; n++ {
		_ = sub.Access()
		//fmt.Println("accessed")
	}
	done <- true
}

// NOTE receiver/pointer dynamics are weird, see https://play.golang.org/p/0Y0nuxEohuP
func publishRandomPubSub(pub *AtomicFloat64Pub, done chan bool) {
	ticker := time.NewTicker(100 * time.Microsecond)
	for {
		select {
		case <-ticker.C:
			pub.Update(rand.Float64())
			// fmt.Println("updated")
		case <-done:
			// entered when done is closed
			return
		}
	}
}

func BenchmarkAtomicFloat64_Transport_Access(b *testing.B) {
	var intensityTransport AtomicFloat64Transport
	done := make(chan bool)

	go publishRandomTransport(&intensityTransport, done)
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

// NOTE receiver/pointer dynamics are weird, see https://play.golang.org/p/0Y0nuxEohuP
func publishRandomTransport(transport *AtomicFloat64Transport, done chan bool) {
	ticker := time.NewTicker(100 * time.Microsecond)
	for {
		select {
		case <-ticker.C:
			transport.Update(rand.Float64())
			// fmt.Println("updated")
		case <-done:
			// entered when done is closed
			return
		}
	}
}

// https://godbolt.org/z/8r8dqz

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
