package transport

import "time"

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ValueType=float64,r3.Vec,quat.Number"

func (avt *AtomicValueTypeTransport) AsPub() *AtomicValueTypePub {
	return &AtomicValueTypePub{avt}
}

func (avt *AtomicValueTypeTransport) AsSub() *AtomicValueTypeSub {
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
