package stride

import "testing"

type Float64Ring []float64
type Vec1Ring []Vec1
type Vec2Ring []Vec2
type VecNRing []VecN

type Vec1 struct {
	X float64
}

type Vec2 struct {
	X float64
	Y float64
}

type VecN struct {
	dim  int
	data []float64
}

func floatToVecN(val float64) VecN {
	return VecN{
		dim:  1,
		data: []float64{val},
	}
}

func BenchmarkBareTypeAlias(b *testing.B) {
	ring := Float64Ring([]float64{
		1.0, 1.1, 1.2, 2.0, 2.1,
	})
	for n := 0; n < b.N; n++ {
		sum := 0.0
		for _, val := range ring {
			sum += val
		}
		_ = sum
	}
}

func BenchmarkSingleElementStruct(b *testing.B) {
	ring := Vec1Ring([]Vec1{
		{X: 1.0}, {X: 1.1}, {X: 1.2}, {X: 2.0}, {X: 2.1},
	})
	for n := 0; n < b.N; n++ {
		sum := 0.0
		for _, val := range ring {
			sum += val.X
		}
		_ = sum
	}
}

func BenchmarkMultiElementStruct(b *testing.B) {
	ring := Vec2Ring([]Vec2{
		{X: 1.0}, {X: 1.1}, {X: 1.2}, {X: 2.0}, {X: 2.1},
	})
	for n := 0; n < b.N; n++ {
		sum := 0.0
		for _, val := range ring {
			sum += val.X
		}
		_ = sum
	}
}

func makeVecNRing() *VecNRing {
	ring := VecNRing([]VecN{
		floatToVecN(1.0),
		floatToVecN(1.1),
		floatToVecN(1.2),
		floatToVecN(2.0),
		floatToVecN(2.1),
	})
	return &ring
}

func BenchmarkDynamicStruct(b *testing.B) {
	data := makeVecNRing()
	for n := 0; n < b.N; n++ {
		sum := 0.0
		for _, val := range *data {
			sum += val.data[0]
		}
		_ = sum
	}
}
