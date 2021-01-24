package accelhist

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
	"testing"
)

type Float64Ring []float64

const elems = 50

func setupRawSlice() *Float64Ring {
	data := Float64Ring(make([]float64, elems))
	for i := range data {
		data[i] = rand.Float64()
	}
	return &data
}

func setupR3Vecs() *[]r3.Vec {
	data := make([]r3.Vec, elems)
	for i := range data {
		data[i] = r3.Vec{
			X: rand.Float64(),
			Y: rand.Float64(),
			Z: rand.Float64(),
		}
	}
	return &data
}

func setupMat() *mat.Dense {
	data := make([]float64, elems*3)
	for i := range data {
		data[i] = rand.Float64()
	}
	return mat.NewDense(3, 50, data)
}

func maxAbs(s []float64) float64 {
	var maxAbs float64
	for _, val := range s {
		if abs := math.Abs(val); abs > maxAbs {
			maxAbs = abs
		}
	}
	return maxAbs
}

func BenchmarkRawSlice_Norm(b *testing.B) {
	rings := [...]*Float64Ring{
		setupRawSlice(), setupRawSlice(), setupRawSlice(),
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, ring := range rings {
			// floats.Norm(s []float64, math.Inf(1)) returns max absolute value in s
			_ = floats.Norm(*ring, math.Inf(1))
		}
	}
}

func BenchmarkRawSlice_MinMax(b *testing.B) {
	rings := [...]*Float64Ring{
		setupRawSlice(), setupRawSlice(), setupRawSlice(),
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, ring := range rings {
			_ = math.Max(-floats.Min(*ring), floats.Max(*ring))
		}
	}
}

func BenchmarkRawSlice_Naive(b *testing.B) {
	rings := [...]*Float64Ring{
		setupRawSlice(), setupRawSlice(), setupRawSlice(),
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, ring := range rings {
			_ = maxAbs(*ring)
		}
	}
}

func BenchmarkR3Vecs_Naive(b *testing.B) {
	vecSlice := setupR3Vecs()
	max := r3.Vec{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, vec := range *vecSlice {
			if abs := math.Abs(vec.X); abs > max.X {
				max.X = abs
			}
			if abs := math.Abs(vec.Y); abs > max.Y {
				max.Y = abs
			}
			if abs := math.Abs(vec.Z); abs > max.Z {
				max.Z = abs
			}
		}
		_ = max
	}
}

func BenchmarkMatrix_MinMax(b *testing.B) {
	matrix := setupMat()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = math.Max(-mat.Min(matrix), mat.Max(matrix))
	}
}

func BenchmarkMatrix_RowViewNorm(b *testing.B) {
	matrix := setupMat()
	b.ResetTimer()
	rows, _ := matrix.Dims()
	for n := 0; n < b.N; n++ {
		for row := 0; row < rows; row++ {
			_ = floats.Norm(matrix.RawRowView(row), math.Inf(1))
		}
	}
}

func BenchmarkMatrix_RowViewMinMax(b *testing.B) {
	matrix := setupMat()
	b.ResetTimer()
	rows, _ := matrix.Dims()
	for n := 0; n < b.N; n++ {
		for row := 0; row < rows; row++ {
			rowView := matrix.RawRowView(row)
			_ = math.Max(-floats.Min(rowView), floats.Max(rowView))
		}
	}
}

func BenchmarkMatrix_RowViewNaive(b *testing.B) {
	matrix := setupMat()
	b.ResetTimer()
	rows, _ := matrix.Dims()
	for n := 0; n < b.N; n++ {
		for row := 0; row < rows; row++ {
			rowView := matrix.RawRowView(row)
			_ = maxAbs(rowView)
		}
	}
}

func TestRowAndColViews(t *testing.T) {
	data := []float64{
		11, 12, 13,
		21, 22, 23,
		31, 32, 33,
	}
	matrix := mat.NewDense(3, 3, data)
	var r, c mat.VecDense
	r.RowViewOf(matrix, 1)
	c.ColViewOf(matrix, 2)
	fmt.Println("row", r)
	fmt.Println("col", c)
}
