package transport

import (
	"math/rand"
	"periph.io/x/conn/v3/physic"
	"testing"
	"time"
)

func BenchmarkFloat64Transports_Access(b *testing.B) {
	benchmarks := []struct {
		name      string
		transport Float64Transport
	}{
		{"Atomic", &AtomicFloat64Transport{}},
		//{"Channel", NewChannelFloat64Transport()},
		//{"Mutex", &MutexFloat64Transport{}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			done := make(chan struct{})

			pub := bm.transport.(Float64Pub)
			sub := bm.transport.(Float64Sub)

			go publishRandom(pub, done)

			if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
				return
			}

			for n := 0; n < b.N; n++ {
				_ = sub.Access()
			}

			close(done)
		})
	}
}

// NOTE receiver/pointer dynamics are weird, see https://play.golang.org/p/0Y0nuxEohuP
func publishRandom(transport Float64Pub, done chan struct{}) {
	ticker := time.NewTicker((physic.KiloHertz).Period())
	for {
		select {
		case <-ticker.C:
			transport.Update(rand.Float64())
			// fmt.Println("updated")
		case _, _ = <-done:
			// entered when done is closed
			return
		}
	}
}
