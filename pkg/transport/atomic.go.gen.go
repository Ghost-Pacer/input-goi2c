// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package transport

import (
	"fmt"
	"sync/atomic"
	"time"

	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

// Lockless float64 transport using sync/atomic.Value.
type AtomicFloat64Transport struct {
	atom atomic.Value // of float64
}

func NewAtomicFloat64() *AtomicFloat64Transport {
	return &AtomicFloat64Transport{}
}

func (avt *AtomicFloat64Transport) EnsureReady(timeout time.Duration, interval time.Duration) error {
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			return fmt.Errorf("%v has passed but atom was never updated", timeout)
		default:
			if avt.atom.Load() != nil {
				return nil
			}
			time.Sleep(interval)
		}
	}
}

func (avt *AtomicFloat64Transport) Access() float64 {
	return avt.atom.Load().(TimedFloat64).value
}

func (avt *AtomicFloat64Transport) AccessTimed() (float64, EventTimings) {
	timedValue := avt.atom.Load().(TimedFloat64)
	timedValue.timings.Accessed = time.Now()
	return timedValue.value, timedValue.timings
}

func (avt *AtomicFloat64Transport) Update(value float64) {
	avt.atom.Store(TimedFloat64{
		value: value,
	})
}

func (avt *AtomicFloat64Transport) UpdateTimed(value float64, sourced time.Time) {
	avt.atom.Store(TimedFloat64{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Lockless r3.Vec transport using sync/atomic.Value.
type AtomicR3VecTransport struct {
	atom atomic.Value // of r3.Vec
}

func NewAtomicR3Vec() *AtomicR3VecTransport {
	return &AtomicR3VecTransport{}
}

func (avt *AtomicR3VecTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			return fmt.Errorf("%v has passed but atom was never updated", timeout)
		default:
			if avt.atom.Load() != nil {
				return nil
			}
			time.Sleep(interval)
		}
	}
}

func (avt *AtomicR3VecTransport) Access() r3.Vec {
	return avt.atom.Load().(TimedR3Vec).value
}

func (avt *AtomicR3VecTransport) AccessTimed() (r3.Vec, EventTimings) {
	timedValue := avt.atom.Load().(TimedR3Vec)
	timedValue.timings.Accessed = time.Now()
	return timedValue.value, timedValue.timings
}

func (avt *AtomicR3VecTransport) Update(value r3.Vec) {
	avt.atom.Store(TimedR3Vec{
		value: value,
	})
}

func (avt *AtomicR3VecTransport) UpdateTimed(value r3.Vec, sourced time.Time) {
	avt.atom.Store(TimedR3Vec{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Lockless quat.Number transport using sync/atomic.Value.
type AtomicQuatNumberTransport struct {
	atom atomic.Value // of quat.Number
}

func NewAtomicQuatNumber() *AtomicQuatNumberTransport {
	return &AtomicQuatNumberTransport{}
}

func (avt *AtomicQuatNumberTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			return fmt.Errorf("%v has passed but atom was never updated", timeout)
		default:
			if avt.atom.Load() != nil {
				return nil
			}
			time.Sleep(interval)
		}
	}
}

func (avt *AtomicQuatNumberTransport) Access() quat.Number {
	return avt.atom.Load().(TimedQuatNumber).value
}

func (avt *AtomicQuatNumberTransport) AccessTimed() (quat.Number, EventTimings) {
	timedValue := avt.atom.Load().(TimedQuatNumber)
	timedValue.timings.Accessed = time.Now()
	return timedValue.value, timedValue.timings
}

func (avt *AtomicQuatNumberTransport) Update(value quat.Number) {
	avt.atom.Store(TimedQuatNumber{
		value: value,
	})
}

func (avt *AtomicQuatNumberTransport) UpdateTimed(value quat.Number, sourced time.Time) {
	avt.atom.Store(TimedQuatNumber{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}