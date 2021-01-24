// Package nearsight provides implementations of a bounded stack data structure.
//
// Right now this is designed to be typed using the interface from client code, providing
// flexibility for the instantiator to use different backing implementations if we decide
// they need to change. There is still the freedom to use the concretely typed backing
// implementations directly. In the future we could change over the package to use a type
// alias for the default implementation if the virtual function call overhead - around 1ns
// or 15% difference @ AMD64 over raw slice access - becomes significant. Then again, if that is
// really an issue it's probably best to just use the backing implementations directly.
//
// A default constructor is provided that currently uses the Ptr
// implementation, since it is only around 15% (1ns) additional overhead in At() to do the
// pointer arithmetic, but Push() has constant runtime with respect to the capacity of the stack.
// The COW (Copy-on-Write) implementation has a slightly faster At() but its Push() must copy
// every element, and therefore it gets quite slow as the size of the stack grows:
// no significant difference with capacity 10, but around 500% slower with capacity 500.
package nearsight

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type ValueType generic.Type

// ValueTypeStack stores a bounded stack of ValueTypes. The capacity is specified at construction time.
// When an element is pushed while the stack is at capacity, the new element acquires index 0; all
// old elements migrate to indices 1, 2, ...; and the oldest element is discarded.
type ValueTypeStack interface {
	// Push adds the element to the front of the stack (index 0) and pushes all existing elements back one index
	Push(value ValueType)
	// At gets the element at the given index [0 newest..Len() oldest)
	At(index int) ValueType

	// Len returns the number of elements stored in the stack (less than or equal to Cap())
	Len() int
	// Cap returns the maximum number of elements that may be stored at once
	Cap() int

	// AsSlice returns the backing slice, for the purpose of iterating the whole stack
	AsSlice() []ValueType
}

// enforce compliance of Ptr implementation
var _ ValueTypeStack = (*PtrValueTypeStack)(nil)

// enforce compliance of COW implementation
var _ ValueTypeStack = (*COWValueTypeStack)(nil)

// NewValueTypeStack returns an initialized ValueTypeStack with the preferred backing implementation
func NewValueTypeStack(capacity int) ValueTypeStack {
	return NewPtrValueTypeStack(capacity)
}
