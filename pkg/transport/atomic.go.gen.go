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

// Lockless float64 transport using sync/atomic.Value. Usable as zero value. Needs pointer semantics.
// Performance note: using dynamic dispatch Pub/Sub interfaces slows down operations by about 450% @ AMD64
// over using static dispatch AtomicPub/AtomicSub, which in turn are about 20% @ AMD64 slower than
// bare references to AtomicTransport.
type AtomicFloat64Transport struct {
	atom atomic.Value // of float64
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
	return avt.atom.Load().(TimedFloat64).Value
}

func (avt *AtomicFloat64Transport) AccessTimed() (float64, EventTimings) {
	timedValue := avt.atom.Load().(TimedFloat64)
	timedValue.Timings.Accessed = time.Now()
	return timedValue.Value, timedValue.Timings
}

func (avt *AtomicFloat64Transport) Update(value float64) {
	avt.atom.Store(TimedFloat64{
		Value: value,
	})
}

func (avt *AtomicFloat64Transport) UpdateTimed(value float64, sourced time.Time) {
	avt.atom.Store(TimedFloat64{
		Value: value,
		Timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Static dispatch Pub/Sub views for higher performance if needed.

func (avt *AtomicFloat64Transport) PubView() *AtomicFloat64Pub {
	return &AtomicFloat64Pub{avt}
}

func (avt *AtomicFloat64Transport) SubView() *AtomicFloat64Sub {
	return &AtomicFloat64Sub{avt}
}

type AtomicFloat64Pub struct {
	transport *AtomicFloat64Transport
}

func (avp *AtomicFloat64Pub) Update(value float64) {
	avp.transport.Update(value)
}

func (avp *AtomicFloat64Pub) UpdateTimed(value float64, sourced time.Time) {
	avp.transport.UpdateTimed(value, sourced)
}

type AtomicFloat64Sub struct {
	transport *AtomicFloat64Transport
}

func (avs *AtomicFloat64Sub) EnsureReady(timeout time.Duration, interval time.Duration) error {
	return avs.transport.EnsureReady(timeout, interval)
}

func (avs *AtomicFloat64Sub) Access() float64 {
	return avs.transport.Access()
}

func (avs *AtomicFloat64Sub) AccessTimed() (float64, EventTimings) {
	return avs.transport.AccessTimed()
}

// Lockless int transport using sync/atomic.Value. Usable as zero value. Needs pointer semantics.
// Performance note: using dynamic dispatch Pub/Sub interfaces slows down operations by about 450% @ AMD64
// over using static dispatch AtomicPub/AtomicSub, which in turn are about 20% @ AMD64 slower than
// bare references to AtomicTransport.
type AtomicIntTransport struct {
	atom atomic.Value // of int
}

func (avt *AtomicIntTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
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

func (avt *AtomicIntTransport) Access() int {
	return avt.atom.Load().(TimedInt).Value
}

func (avt *AtomicIntTransport) AccessTimed() (int, EventTimings) {
	timedValue := avt.atom.Load().(TimedInt)
	timedValue.Timings.Accessed = time.Now()
	return timedValue.Value, timedValue.Timings
}

func (avt *AtomicIntTransport) Update(value int) {
	avt.atom.Store(TimedInt{
		Value: value,
	})
}

func (avt *AtomicIntTransport) UpdateTimed(value int, sourced time.Time) {
	avt.atom.Store(TimedInt{
		Value: value,
		Timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Static dispatch Pub/Sub views for higher performance if needed.

func (avt *AtomicIntTransport) PubView() *AtomicIntPub {
	return &AtomicIntPub{avt}
}

func (avt *AtomicIntTransport) SubView() *AtomicIntSub {
	return &AtomicIntSub{avt}
}

type AtomicIntPub struct {
	transport *AtomicIntTransport
}

func (avp *AtomicIntPub) Update(value int) {
	avp.transport.Update(value)
}

func (avp *AtomicIntPub) UpdateTimed(value int, sourced time.Time) {
	avp.transport.UpdateTimed(value, sourced)
}

type AtomicIntSub struct {
	transport *AtomicIntTransport
}

func (avs *AtomicIntSub) EnsureReady(timeout time.Duration, interval time.Duration) error {
	return avs.transport.EnsureReady(timeout, interval)
}

func (avs *AtomicIntSub) Access() int {
	return avs.transport.Access()
}

func (avs *AtomicIntSub) AccessTimed() (int, EventTimings) {
	return avs.transport.AccessTimed()
}

// Lockless r3.Vec transport using sync/atomic.Value. Usable as zero value. Needs pointer semantics.
// Performance note: using dynamic dispatch Pub/Sub interfaces slows down operations by about 450% @ AMD64
// over using static dispatch AtomicPub/AtomicSub, which in turn are about 20% @ AMD64 slower than
// bare references to AtomicTransport.
type AtomicR3VecTransport struct {
	atom atomic.Value // of r3.Vec
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
	return avt.atom.Load().(TimedR3Vec).Value
}

func (avt *AtomicR3VecTransport) AccessTimed() (r3.Vec, EventTimings) {
	timedValue := avt.atom.Load().(TimedR3Vec)
	timedValue.Timings.Accessed = time.Now()
	return timedValue.Value, timedValue.Timings
}

func (avt *AtomicR3VecTransport) Update(value r3.Vec) {
	avt.atom.Store(TimedR3Vec{
		Value: value,
	})
}

func (avt *AtomicR3VecTransport) UpdateTimed(value r3.Vec, sourced time.Time) {
	avt.atom.Store(TimedR3Vec{
		Value: value,
		Timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Static dispatch Pub/Sub views for higher performance if needed.

func (avt *AtomicR3VecTransport) PubView() *AtomicR3VecPub {
	return &AtomicR3VecPub{avt}
}

func (avt *AtomicR3VecTransport) SubView() *AtomicR3VecSub {
	return &AtomicR3VecSub{avt}
}

type AtomicR3VecPub struct {
	transport *AtomicR3VecTransport
}

func (avp *AtomicR3VecPub) Update(value r3.Vec) {
	avp.transport.Update(value)
}

func (avp *AtomicR3VecPub) UpdateTimed(value r3.Vec, sourced time.Time) {
	avp.transport.UpdateTimed(value, sourced)
}

type AtomicR3VecSub struct {
	transport *AtomicR3VecTransport
}

func (avs *AtomicR3VecSub) EnsureReady(timeout time.Duration, interval time.Duration) error {
	return avs.transport.EnsureReady(timeout, interval)
}

func (avs *AtomicR3VecSub) Access() r3.Vec {
	return avs.transport.Access()
}

func (avs *AtomicR3VecSub) AccessTimed() (r3.Vec, EventTimings) {
	return avs.transport.AccessTimed()
}

// Lockless quat.Number transport using sync/atomic.Value. Usable as zero value. Needs pointer semantics.
// Performance note: using dynamic dispatch Pub/Sub interfaces slows down operations by about 450% @ AMD64
// over using static dispatch AtomicPub/AtomicSub, which in turn are about 20% @ AMD64 slower than
// bare references to AtomicTransport.
type AtomicQuatNumberTransport struct {
	atom atomic.Value // of quat.Number
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
	return avt.atom.Load().(TimedQuatNumber).Value
}

func (avt *AtomicQuatNumberTransport) AccessTimed() (quat.Number, EventTimings) {
	timedValue := avt.atom.Load().(TimedQuatNumber)
	timedValue.Timings.Accessed = time.Now()
	return timedValue.Value, timedValue.Timings
}

func (avt *AtomicQuatNumberTransport) Update(value quat.Number) {
	avt.atom.Store(TimedQuatNumber{
		Value: value,
	})
}

func (avt *AtomicQuatNumberTransport) UpdateTimed(value quat.Number, sourced time.Time) {
	avt.atom.Store(TimedQuatNumber{
		Value: value,
		Timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Static dispatch Pub/Sub views for higher performance if needed.

func (avt *AtomicQuatNumberTransport) PubView() *AtomicQuatNumberPub {
	return &AtomicQuatNumberPub{avt}
}

func (avt *AtomicQuatNumberTransport) SubView() *AtomicQuatNumberSub {
	return &AtomicQuatNumberSub{avt}
}

type AtomicQuatNumberPub struct {
	transport *AtomicQuatNumberTransport
}

func (avp *AtomicQuatNumberPub) Update(value quat.Number) {
	avp.transport.Update(value)
}

func (avp *AtomicQuatNumberPub) UpdateTimed(value quat.Number, sourced time.Time) {
	avp.transport.UpdateTimed(value, sourced)
}

type AtomicQuatNumberSub struct {
	transport *AtomicQuatNumberTransport
}

func (avs *AtomicQuatNumberSub) EnsureReady(timeout time.Duration, interval time.Duration) error {
	return avs.transport.EnsureReady(timeout, interval)
}

func (avs *AtomicQuatNumberSub) Access() quat.Number {
	return avs.transport.Access()
}

func (avs *AtomicQuatNumberSub) AccessTimed() (quat.Number, EventTimings) {
	return avs.transport.AccessTimed()
}
