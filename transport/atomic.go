package transport

import (
	"fmt"
	"sync/atomic"
	"time"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ValueType=float64,r3.Vec,quat.Number"

// Lockless ValueType transport using sync/atomic.Value.
type AtomicValueTypeTransport struct {
	atom atomic.Value // of ValueType
}

func (avt *AtomicValueTypeTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
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

func (avt *AtomicValueTypeTransport) Access() ValueType {
	return avt.atom.Load().(TimedValueType).value
}

func (avt *AtomicValueTypeTransport) AccessTimed() (ValueType, EventTimings) {
	timedValue := avt.atom.Load().(TimedValueType)
	timedValue.timings.Accessed = time.Now()
	return timedValue.value, timedValue.timings
}

func (avt *AtomicValueTypeTransport) Update(value ValueType) {
	avt.atom.Store(TimedValueType{
		value: value,
	})
}

func (avt *AtomicValueTypeTransport) UpdateTimed(value ValueType, sourced time.Time) {
	avt.atom.Store(TimedValueType{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}
