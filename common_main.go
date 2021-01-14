package main

import (
	"fmt"
	"github.com/Ghost-Pacer/input-goi2c/transport"
	"math/rand"
	"os"
	"time"
)

func BenchmarkFloat64Transports_Access() {
	done := make(chan bool)
	var trans transport.AtomicFloat64Transport
	var sub = &trans

	go publishRandom(&trans, done)

	if err := sub.EnsureReady(time.Second, time.Millisecond); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	for n := 0; n < 10000000; n++ {
		_ = sub.Access()
	}

	done <- true
}

// NOTE receiver/pointer dynamics are weird, see https://play.golang.org/p/0Y0nuxEohuP
func publishRandom(transport *transport.AtomicFloat64Transport, done chan bool) {
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

func main() {
	BenchmarkFloat64Transports_Access()
}
