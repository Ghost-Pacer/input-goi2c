package transport

import (
	"fmt"
	"math/rand"
	"os"
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
			done := make(chan bool)
			var trans AtomicFloat64Transport
			var sub *AtomicFloat64Transport = &trans

			go publishRandom(sub, done)

			if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
				return
			}

			func(sub Float64Sub) {
				for n := 0; n < b.N; n++ {
					_ = sub.Access()
				}
			}(sub)

			done <- true
		})
	}
}

// NOTE receiver/pointer dynamics are weird, see https://play.golang.org/p/0Y0nuxEohuP
func publishRandom(transport *AtomicFloat64Transport, done chan bool) {
	ticker := time.NewTicker((10 * physic.KiloHertz).Period())
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
