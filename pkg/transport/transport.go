// Package transport implements high-performance abstractions for sharing data between goroutines.
package transport

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=float64,int,r3.Vec,quat.Number"

import (
	"time"

	"github.com/cheekybits/genny/generic"
)

type ValueType generic.Type

type TimedValueType struct {
	Value   ValueType
	Timings EventTimings
}

// Interface for publishers; analogue of chan<- ValueType
type ValueTypePub interface {
	Update(value ValueType)
	UpdateTimed(value ValueType, source time.Time)
}

// Interface for subscribers; analogue of <-chan ValueType
type ValueTypeSub interface {
	EnsureReady(interval time.Duration, timeout time.Duration) error
	Access() ValueType
	AccessTimed() (ValueType, EventTimings)
}

var _ ValueTypePub = (*AtomicValueTypeTransport)(nil)
var _ ValueTypeSub = (*AtomicValueTypeTransport)(nil)
