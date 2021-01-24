package nearsight

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type PtrValueTypeStack struct {
	data      []ValueType
	headIndex int
}

// Implemented using pointer semantics because we need to be able to modify the header of the data slice
// when things are first being pushed in. This way, len(data) actually returns an accurate value when
// the stack is not full. We could preallocate having len == cap and then we could get value semantics
// for the whole collection, but then we would wind up with zero values in the slice when the stack is
// not full. The stuff that goes in here should (hopefully) be just value types
func NewPtrValueTypeStack(capacity int) *PtrValueTypeStack {
	s := PtrValueTypeStack{
		data:      make([]ValueType, 0, capacity),
		headIndex: -1,
	}
	return &s
}

func (s *PtrValueTypeStack) Push(value ValueType) {
	// increment the head index and wrap around
	s.headIndex = (s.headIndex + 1) % cap(s.data)
	switch {
	case len(s.data) <= s.headIndex:
		// if there aren't already cap elements, append must be used to extend the length
		s.data = append(s.data, value)
	default:
		s.data[s.headIndex] = value
	}
}

func (s *PtrValueTypeStack) At(index int) ValueType {
	// since the head index moves forward with each new element, we have to look backwards to get newest things first
	// also, add an additional cap before modulo in case (-index + s.headIndex) needs to wrap around
	// because modulo will not automatically wrap up to a positive value
	realIndex := (-index + s.headIndex + cap(s.data)) % cap(s.data)
	return s.data[realIndex]
}

func (s *PtrValueTypeStack) Cap() int {
	return cap(s.data)
}

func (s *PtrValueTypeStack) Len() int {
	return len(s.data)
}

func (s *PtrValueTypeStack) AsSlice() []ValueType {
	return s.data
}
