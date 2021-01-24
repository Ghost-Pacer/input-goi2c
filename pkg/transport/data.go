// Package transport implements high-performance abstractions for sharing data between goroutines.
package transport

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=float64,r3.Vec,quat.Number"

import (
	"github.com/cheekybits/genny/generic"
)

type ValueType generic.Type

type TimedValueType struct {
	value   ValueType
	timings EventTimings
}
