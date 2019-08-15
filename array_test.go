package golina

import (
	"log"
	"math/rand"
	"os"
	"runtime/trace"
	"testing"
	"time"
)

func generateRandomVector(size int) *Vector {
	slice := make(Vector, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Float64() - rand.Float64()
	}
	return &slice
}

// https://blog.karenuorteva.fi/go-unit-test-setup-and-teardown-db1601a796f2#.2aherx2z5

func TestMatrix_Init(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA._array == nil {
		t.Fail()
	}
}

func TestMatrix_Dims(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row, col := matA.Dims()
	if row != 3 || col != 3 {
		t.Fail()
	}
}

func TestMatrix_At(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA.At(1, 1) != 5 {
		t.Fail()
	}
}

func TestMatrix_Set(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matA.Set(1, 1, 10)
	if matA.At(1, 1) != 10 {
		t.Fail()
	}
}

func TestMatrix_T(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matAT := matA.T()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matAT.At(i, j) != matA.At(j, i) {
				t.Fail()
			}
		}
	}
}

func TestConvolve(t *testing.T) {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()
	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()
	size := 10000
	u := generateRandomVector(size)
	v := generateRandomVector(size)

	res := Convolve(u, v)
	if len(*res) != size+size-1 {
		t.Fail()
	}
}
