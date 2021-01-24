package cow

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type ValueType generic.Type

type ValueTypeStack []ValueType

func NewValueTypeStack(size int) *ValueTypeStack {
	s := ValueTypeStack(make([]ValueType, 0, size))
	return &s
}

func (s *ValueTypeStack) Push(value ValueType) {
	// store of append needs pointer semantics to modify slice header
	// so I used pointer semantics everywhere to avoid two symbols for the same thing
	// but pointer semantics are ugly idk
	switch {
	case len(*s) == 0:
		*s = append(*s, value)
	case len(*s) < cap(*s):
		*s = append(*s, value)
		copy((*s)[1:], *s)
		(*s)[0] = value
	default:
		copy((*s)[1:], *s)
		(*s)[0] = value
	}
}

func (s *ValueTypeStack) At(index int) ValueType {
	return (*s)[index]
}

func (s *ValueTypeStack) Size() int {
	return cap(*s)
}
