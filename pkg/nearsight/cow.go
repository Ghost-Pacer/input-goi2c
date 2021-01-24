package nearsight

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type COWValueTypeStack []ValueType

func NewCOWValueTypeStack(size int) *COWValueTypeStack {
	s := COWValueTypeStack(make([]ValueType, 0, size))
	return &s
}

func (s *COWValueTypeStack) Push(value ValueType) {
	// store of append needs pointer semantics to modify slice header
	// so I used pointer semantics everywhere to avoid two symbols for the same thing
	// but pointer semantics are ugly idk
	switch {
	case len(*s) == 0:
		// zero elements needs a special case because indexing (*s)[1:] will fail if *s is empty
		*s = append(*s, value)
	case len(*s) < cap(*s):
		// need to extend the length by appending something
		*s = append(*s, value)
		// this will copy over the previous append but now the slice is long enough
		copy((*s)[1:], *s)
		// replace the value in the correct location
		(*s)[0] = value
	default:
		// stack is already full, just copy over the elements and then insert
		copy((*s)[1:], *s)
		(*s)[0] = value
	}
}

func (s *COWValueTypeStack) At(index int) ValueType {
	return (*s)[index]
}

func (s *COWValueTypeStack) Cap() int {
	return cap(*s)
}

func (s *COWValueTypeStack) Len() int {
	return len(*s)
}

func (s *COWValueTypeStack) AsSlice() []ValueType {
	return *s
}
