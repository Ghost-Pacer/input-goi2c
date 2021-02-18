package transport

import (
	"fmt"
	"sync/atomic"
	"time"
)

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=float64,int,r3.Vec,quat.Number"

// Lockless ValueType transport using sync/atomic.Value. Usable as zero value. Needs pointer semantics.
// Performance note: using dynamic dispatch Pub/Sub interfaces slows down operations by about 450% @ AMD64
// over using static dispatch AtomicPub/AtomicSub, which in turn are about 20% @ AMD64 slower than
// bare references to AtomicTransport.
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
	return avt.atom.Load().(TimedValueType).Value
}

func (avt *AtomicValueTypeTransport) AccessTimed() (ValueType, EventTimings) {
	timedValue := avt.atom.Load().(TimedValueType)
	timedValue.Timings.Accessed = time.Now()
	return timedValue.Value, timedValue.Timings
}

func (avt *AtomicValueTypeTransport) Update(value ValueType) {
	avt.atom.Store(TimedValueType{
		Value: value,
	})
}

func (avt *AtomicValueTypeTransport) UpdateTimed(value ValueType, sourced time.Time) {
	avt.atom.Store(TimedValueType{
		Value: value,
		Timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	})
}

// Static dispatch Pub/Sub views for higher performance if needed.

func (avt *AtomicValueTypeTransport) PubView() *AtomicValueTypePub {
	return &AtomicValueTypePub{avt}
}

func (avt *AtomicValueTypeTransport) SubView() *AtomicValueTypeSub {
	return &AtomicValueTypeSub{avt}
}

type AtomicValueTypePub struct {
	transport *AtomicValueTypeTransport
}

func (avp *AtomicValueTypePub) Update(value ValueType) {
	avp.transport.Update(value)
}

func (avp *AtomicValueTypePub) UpdateTimed(value ValueType, sourced time.Time) {
	avp.transport.UpdateTimed(value, sourced)
}

type AtomicValueTypeSub struct {
	transport *AtomicValueTypeTransport
}

func (avs *AtomicValueTypeSub) EnsureReady(timeout time.Duration, interval time.Duration) error {
	return avs.transport.EnsureReady(timeout, interval)
}

func (avs *AtomicValueTypeSub) Access() ValueType {
	return avs.transport.Access()
}

func (avs *AtomicValueTypeSub) AccessTimed() (ValueType, EventTimings) {
	return avs.transport.AccessTimed()
}
