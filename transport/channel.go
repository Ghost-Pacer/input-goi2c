package transport

import (
	"fmt"
	"time"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ValueType=float64,r3.Vec,quat.Number"

// Channel-based ValueType transport with most-recent-value caching & synchronous calling convention
type ChannelValueTypeTransport struct {
	channel chan TimedValueType
	cache   TimedValueType
}

func NewChannelValueTypeTransport() *ChannelValueTypeTransport {
	return &ChannelValueTypeTransport{
		channel: make(chan TimedValueType, 1),
	}
}

func (cvt *ChannelValueTypeTransport) EnsureReady(timeout time.Duration, interval time.Duration) error {
	select {
	case update := <-cvt.channel:
		cvt.cache = update
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("%v has passed but channel never received", timeout)
	}
}

// FIXME: if the consumer is relatively slow, the value returned will be the first value produced
//  after the previous Access() call.
func (cvt *ChannelValueTypeTransport) Access() ValueType {
	select {
	case update := <-cvt.channel:
		cvt.cache = update
		return update.value
	default:
		return cvt.cache.value
	}
}

// FIXME: if the consumer is relatively slow, the value returned will be the first value produced
//  after the previous AccessTimed() call.
func (cvt *ChannelValueTypeTransport) AccessTimed() (ValueType, EventTimings) {
	select {
	case update := <-cvt.channel:
		update.timings.Accessed = time.Now()
		cvt.cache = update
		return update.value, update.timings
	default:
		return cvt.cache.value, cvt.cache.timings
	}
}

func (cvt *ChannelValueTypeTransport) Update(value ValueType) {
	select {
	case cvt.channel <- TimedValueType{value: value}:
		// write succeeded, do nothing
	}
	// otherwise do nothing
}

func (cvt *ChannelValueTypeTransport) UpdateTimed(value ValueType, sourced time.Time) {
	update := TimedValueType{
		value: value,
		timings: EventTimings{
			Sourced: sourced,
			Updated: time.Now(),
		},
	}
	select {
	case cvt.channel <- update:
		// write succeeded, do nothing
	}
	// otherwise do nothing
}
