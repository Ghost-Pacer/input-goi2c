package nearsight

import (
	"github.com/Ghost-Pacer/input-goi2c/nearsight/ptr"
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=$GOFILE.gen.go gen "ValueType=mat.VecDense,r3.Vec,float64"

type ValueType generic.Type

type ValueTypeStack = ptr.ValueTypeStack

func NewValueTypes(size int) *ValueTypeStack {
	return ptr.NewValueTypeStack(size)
}
