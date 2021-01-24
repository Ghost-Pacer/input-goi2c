package cow

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,float64"

type ValueType generic.Type

type ValueTypeRing []ValueType

func NewValueTypeRing(size int) *ValueTypeRing {
	r := ValueTypeRing(make([]ValueType, 0, size))
	return &r
}

func (r *ValueTypeRing) Push(value ValueType) {
	// deref for convenience, slice header does not change so direct write to pointer not needed
	sl := *r
	// copy elements 1..last to positions 0..last-1
	copy(sl, sl[1:])
	// insert value at end of slice
	sl[cap(sl)] = value
}

func (r *ValueTypeRing) At(index int) ValueType {
	if index < 0 {
		// support negative indexing
		index += len(*r)
	}
	return (*r)[index]
}
