package nearsight

import (
	"math/rand"
	"testing"
)

func BenchmarkBarePtrPush10(b *testing.B) {
	queue := NewFloat64Stack(10)
	for n := 0; n < b.N; n++ {
		queue.Push(rand.Float64())
	}
}

func BenchmarkBarePtrPush500(b *testing.B) {
	queue := NewPtrFloat64Stack(500)
	for n := 0; n < b.N; n++ {
		queue.Push(rand.Float64())
	}
}

func BenchmarkBarePtrAt10(b *testing.B) {
	queue := NewPtrFloat64Stack(10)
	for i := 0; i < 20; i++ {
		queue.Push(rand.Float64())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = queue.At(rand.Intn(10))
	}
}

func BenchmarkBarePtrAt500(b *testing.B) {
	queue := NewPtrFloat64Stack(500)
	for i := 0; i < 1000; i++ {
		queue.Push(rand.Float64())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = queue.At(rand.Intn(500))
	}
}

func BenchmarkBareCowPush10(b *testing.B) {
	queue := NewCOWFloat64Stack(10)
	for n := 0; n < b.N; n++ {
		queue.Push(rand.Float64())
	}
}

func BenchmarkBareCowPush500(b *testing.B) {
	queue := NewCOWFloat64Stack(500)
	for n := 0; n < b.N; n++ {
		queue.Push(rand.Float64())
	}
}

func BenchmarkBareCowAt10(b *testing.B) {
	queue := NewCOWFloat64Stack(10)
	for i := 0; i < 20; i++ {
		queue.Push(rand.Float64())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = queue.At(rand.Intn(10))
	}
}

func BenchmarkBareCowAt500(b *testing.B) {
	queue := NewCOWFloat64Stack(500)
	for i := 0; i < 1000; i++ {
		queue.Push(rand.Float64())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = queue.At(rand.Intn(500))
	}
}

func ifacePush(b *testing.B, queue Float64Stack) {
	for n := 0; n < b.N; n++ {
		queue.Push(rand.Float64())
	}
}

func ifaceAt(b *testing.B, queue Float64Stack) {
	size := queue.Cap()
	for i := 0; i < size*2; i++ {
		queue.Push(rand.Float64())
	}
	for n := 0; n < b.N; n++ {
		_ = queue.At(rand.Intn(size))
	}
}

func BenchmarkIfacePtrPush10(b *testing.B) {
	ifacePush(b, NewPtrFloat64Stack(10))
}

func BenchmarkIfacePtrPush500(b *testing.B) {
	ifacePush(b, NewPtrFloat64Stack(500))
}

func BenchmarkIfaceCowPush10(b *testing.B) {
	ifacePush(b, NewCOWFloat64Stack(10))
}

func BenchmarkIfaceCowPush500(b *testing.B) {
	ifacePush(b, NewCOWFloat64Stack(500))
}

func BenchmarkIfacePtrAt10(b *testing.B) {
	ifaceAt(b, NewPtrFloat64Stack(10))
}

func BenchmarkIfacePtrAt500(b *testing.B) {
	ifaceAt(b, NewPtrFloat64Stack(500))
}

func BenchmarkIfaceCowAt10(b *testing.B) {
	ifaceAt(b, NewCOWFloat64Stack(10))
}

func BenchmarkIfaceCowAt500(b *testing.B) {
	ifaceAt(b, NewCOWFloat64Stack(500))
}
