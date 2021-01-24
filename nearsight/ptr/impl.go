package ptr

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type ValueType generic.Type

type ValueTypeStack struct {
	data      []ValueType
	headIndex int
}

func NewValueTypeStack(size int) *ValueTypeStack {
	s := ValueTypeStack{
		data:      make([]ValueType, 0, size),
		headIndex: -1,
	}
	return &s
}

func (s *ValueTypeStack) Push(value ValueType) {
	s.headIndex = (s.headIndex + 1) % cap(s.data)
	switch {
	case len(s.data) <= s.headIndex:
		s.data = append(s.data, value)
	default:
		s.data[s.headIndex] = value
	}
}

func (s *ValueTypeStack) At(index int) ValueType {
	realIndex := (-index + s.headIndex + cap(s.data)) % cap(s.data)
	return s.data[realIndex]
}

func (s *ValueTypeStack) Size() int {
	return cap(s.data)
}
