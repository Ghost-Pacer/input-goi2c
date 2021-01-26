package transport

import (
	"math/rand"
	"testing"
	"time"
)

// https://godbolt.org/z/8r8dqz for messing around with static vs dynamic dispatch
// https://play.golang.org/p/0Y0nuxEohuP for weird pointer-receiver interactions

func BenchmarkAtomicFloat64_Access_Dynamic(b *testing.B) {
	transport := new(AtomicFloat64Transport)
	var pub Float64Pub = transport
	var sub Float64Sub = transport

	done := make(chan struct{})
	defer close(done)

	go func() {
		ticker := time.NewTicker(100 * time.Microsecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				pub.Update(rand.Float64())
			}
		}
	}()
	if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
		b.Fatal(err)
		return
	}

	for n := 0; n < b.N; n++ {
		_ = sub.Access()
	}
}

func BenchmarkAtomicFloat64_Access_Static(b *testing.B) {
	transport := new(AtomicFloat64Transport)
	pub := transport.PubView()
	sub := transport.SubView()

	done := make(chan struct{})
	defer close(done)

	go func() {
		ticker := time.NewTicker(100 * time.Microsecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				pub.Update(rand.Float64())
			}
		}
	}()

	if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
		b.Fatal(err)
		return
	}

	for n := 0; n < b.N; n++ {
		_ = sub.Access()
	}
}

func BenchmarkAtomicFloat64_Access_Bare(b *testing.B) {
	transport := new(AtomicFloat64Transport)

	done := make(chan struct{})
	defer close(done)

	go func() {
		ticker := time.NewTicker(100 * time.Microsecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				transport.Update(rand.Float64())
			}
		}
	}()

	if err := transport.EnsureReady(time.Second, time.Millisecond); err != nil {
		b.Fatal(err)
		return
	}

	for n := 0; n < b.N; n++ {
		_ = transport.Access()
	}
}

func TestAtomicGated(t *testing.T) {
	tests := []struct {
		name        string
		cycleWrites int
		cycleReads  int
	}{
		{"Interleaved", 1, 1},
		{"Slow Publisher", 1, 5},
		{"Slow Subscriber", 5, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := new(AtomicIntTransport)
			gate := make(chan struct{})
			defer close(gate)

			// set up publisher
			go func() {
				var pubCount int
				for range gate {
					// send cycleWrites new values every time the gate is released
					for i := 0; i < test.cycleWrites; i++ {
						transport.Update(pubCount)
						pubCount++
					}
				}
			}()

			// prime and verify EnsureReady
			gate <- struct{}{}
			if err := transport.EnsureReady(time.Second, time.Millisecond); err != nil {
				t.Fatal(err)
				return
			}

			for cycle := 0; cycle < 100; cycle++ {
				targetValue := cycle*test.cycleWrites + (test.cycleWrites - 1)
				for i := 0; i < test.cycleReads; i++ {
					if v := transport.Access(); v != targetValue {
						t.Errorf("expected publisher to be at %v but actually got %v", targetValue, v)
					}
				}

				gate <- struct{}{}
				// give time for publisher to receive and act on gate release
				time.Sleep(10 * time.Microsecond)
			}
		})
	}
}
