package transport

import (
	"errors"
	"sync"
	"time"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ValueType=float64,r3.Vec,quat.Number"

/*
ValueType transport using sync.RWMutex.
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
type MutexValueTypeTransport struct {
	mutex   sync.RWMutex
	value   ValueType
	timings EventTimings
	ready   chan bool
}

func (mvt *MutexValueTypeTransport) EnsureReady(timeout time.Duration) error {
	return errors.New("cannot EnsureReady on a MutexValueTypeTransport!")
}

func (mvt *MutexValueTypeTransport) Access() ValueType {
	mvt.mutex.RLock()

	value := mvt.value
	mvt.mutex.RUnlock()
	return value
}

func (mvt *MutexValueTypeTransport) AccessTimed() (ValueType, EventTimings) {
	mvt.mutex.RLock()

	mvt.timings.Accessed = time.Now()

	value, timings := mvt.value, mvt.timings
	mvt.mutex.RUnlock()
	return value, timings
}

func (mvt *MutexValueTypeTransport) Update(value ValueType) {
	mvt.mutex.Lock()

	mvt.value = value

	mvt.mutex.Unlock()
}

func (mvt *MutexValueTypeTransport) UpdateTimed(value ValueType, sourced time.Time) {
	mvt.mutex.Lock()

	mvt.timings.Sourced = sourced
	mvt.timings.Updated = time.Now()
	mvt.value = value

	mvt.mutex.Unlock()
}
