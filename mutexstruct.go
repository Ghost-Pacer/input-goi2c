package main

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
	"gonum.org/v1/gonum/num/quat"
	"sync"
	"sync/atomic"
	"time"
)

type ValueType generic.Type

type EventTimings struct {
	Sourced  time.Time
	Updated  time.Time
	Accessed time.Time
}

type TimedValueType struct {
	value   ValueType
	timings EventTimings
}

type ValueSub interface {
	EnsureReady(timeout time.Duration, interval time.Duration) error
	Access() ValueType
	AccessTimed() (ValueType, EventTimings)
}

type ValuePub interface {
	Update(ValueType)
	UpdateTimed(ValueType, time.Time)
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
	value   ValueType
	timings EventTimings
	ready   chan bool
}

func (mvt *MutexedValueTransport) EnsureReady(timeout time.Duration) error {
	panic("cannot ensure ready on MutexedValueTransport")
}

func (mvt *MutexedValueTransport) Access() ValueType {
	mvt.mutex.RLock()

	value := mvt.value
	mvt.mutex.RUnlock()
	return value
}

func (mvt *MutexedValueTransport) AccessTimed() (ValueType, EventTimings) {
	mvt.mutex.RLock()

	mvt.timings.Accessed = time.Now()

	value, timings := mvt.value, mvt.timings
	mvt.mutex.RUnlock()
	return value, timings
}

func (mvt *MutexedValueTransport) Update(value ValueType) {
	mvt.mutex.Lock()

	mvt.value = value

	mvt.mutex.Unlock()
}

func (mvt *MutexedValueTransport) UpdateTimed(value ValueType, sourced time.Time) {
	mvt.mutex.Lock()

	mvt.timings.Sourced = sourced
	mvt.timings.Updated = time.Now()
	mvt.value = value

	mvt.mutex.Unlock()
}

type AtomicValueTransport struct {
	atom atomic.Value // of ValueType
}

func (avt *AtomicValueTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			return fmt.Errorf("%v has passed but atom was never updated")
		default:
			if avt.atom.Load() != nil {
				return nil
			}
			time.Sleep(interval)
		}
	}
}

func (avt *AtomicValueTransport) Access() ValueType {
	return avt.atom.Load().(TimedValueType).value
}

func (avt *AtomicValueTransport) AccessTimed() (ValueType, EventTimings) {
	timedValue := avt.atom.Load().(TimedValueType)
	timedValue.timings.Accessed = time.Now()
	return timedValue.value, timedValue.timings
}

func (avt *AtomicValueTransport) Update(value ValueType) {
	avt.atom.Store(TimedValueType{
		value: value,
	})
}

func (avt *AtomicValueTransport) UpdateTimed(value ValueType, sourced time.Time) {
	avt.atom.Store(TimedValueType{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

type ChanneledValueTransport struct {
	channel chan TimedValueType
	cache   TimedValueType
	timeout time.Duration
}

func (cvt *ChanneledValueTransport) New(timeout time.Duration) {

}

func (cvt *ChanneledValueTransport) Access() ValueType {

}
