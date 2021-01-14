// Package transport implements high-performance abstractions for sharing data between goroutines.
package transport

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "ValueType=float64,r3.Vec,quat.Number"

import (
	"github.com/cheekybits/genny/generic"
	"time"
)

type ValueType generic.Type

type TimedValueType struct {
	value   ValueType
	timings EventTimings
}

type ValueTypeSub interface {
	EnsureReady(timeout time.Duration, interval time.Duration) error
	Access() ValueType
	AccessTimed() (ValueType, EventTimings)
}

type ValueTypePub interface {
	Update(ValueType)
	UpdateTimed(ValueType, time.Time)
}

type ValueTypeTransport interface {
	ValueTypePub
	ValueTypeSub
}
