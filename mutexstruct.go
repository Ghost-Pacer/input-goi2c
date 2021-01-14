package main

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
	"gonum.org/v1/gonum/spatial/r3"
	"sync"
	"time"
)

type Value generic.Type

type EventTimings struct {
	Sourced  time.Time
	Updated  time.Time
	Accessed time.Time
}

type TimedValue struct {
	value   Value
	timings EventTimings
}

type ValueSub interface {
	EnsureReady(timeout time.Duration) error
	Access() Value
	AccessTimed() (Value, EventTimings)
}

type ValuePub interface {
	Update(Value)
	UpdateTimed(Value, time.Time)
}

/*
High-performance transport for producer-consumer pattern.
Mutex-based transport receivers are traditionally written as
	mutex.Lock()
	defer mutex.Unlock()
	...

However, defer adds a ~16x overhead over simply
	mutex.Lock()
	...
	mutex.Unlock()

and since this transport is designed to go into the production-configuration
hotpath, all of its receivers use the latter pattern. See:
https://medium.com/i0exception/runtime-overhead-of-using-defer-in-go-7140d5c40e32
*/
type MutexedValueTransport struct {
	mutex   sync.RWMutex
	value   Value
	timings EventTimings
	ready chan bool
}

func (mvt *MutexedValueTransport) EnsureReady(timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return fmt.Errorf("%v has passed but MutexedValueTransport was not ready", timeout)
		case
	}
}

func (mvt *MutexedValueTransport) Access() Value {
	mvt.mutex.RLock()

	value := mvt.value
	mvt.mutex.RUnlock()
	return value
}

func (mvt *MutexedValueTransport) AccessTimed() (Value, EventTimings) {
	mvt.mutex.RLock()

	mvt.timings.Accessed = time.Now()

	value, timings := mvt.value, mvt.timings
	mvt.mutex.RUnlock()
	return value, timings
}

func (mvt *MutexedValueTransport) Update(value Value) {
	mvt.mutex.Lock()

	mvt.value = value

	mvt.mutex.Unlock()
}

func (mvt *MutexedValueTransport) UpdateTimed(value Value, sourced time.Time) {
	mvt.mutex.Lock()

	mvt.timings.Sourced = sourced
	mvt.timings.Updated = time.Now()
	mvt.value = value

	mvt.mutex.Unlock()
}

type ChanneledValueTransport struct {
	channel chan TimedValue
	cache TimedValue
	timeout time.Duration
}

func (cvt *ChanneledValueTransport) New(timeout time.Duration) {

}

func (cvt *ChanneledValueTransport) Access() Value {

}
