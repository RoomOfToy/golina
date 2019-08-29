package golina

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestEigenDecompose(t *testing.T) {
	matA := GenerateRandomSymmetric33Matrix()
	V, D := EigenDecompose(matA)
	eval := EigenValues33(matA)
	d := ZeroMatrix(3, 3)
	for i, e := range *eval {
		d.Set(i, i, e)
	}
	evec := EigenVector33(matA, eval)
	if !MEqual(D, d) {
		t.Fail()
	}
	for i := range V._array {
		if !VEqual(V.Col(i), evec.Row(i)) && !VEqual(V.Col(i), evec.Row(i).MulNum(-1.)) {
			t.Fail()
		}
	}

	// non-symmetric
	// EigenDecompose result not in order
	matB := new(Matrix).Init(Data{{1, 2, 3}, {7, 5, 6}, {7, 4, 9}})
	VB, DB := EigenDecompose(matB)
	if !MEqual(matB, VB.Mul(DB).Mul(VB.Inverse())) {
		t.Fail()
	}
}

func TestEigenValues33(t *testing.T) {
	a := Data{{1, 3, 4}, {3, 2, 7}, {4, 7, 5}}
	matA := new(Matrix).Init(a)
	if !VEqual(EigenValues33(matA), &Vector{-3.67018839, -1.10871847, 12.77890686}) {
		t.Fail()
	}
}

func TestEigenVector33(t *testing.T) {
	a := Data{{1, 3, 4}, {3, 2, 7}, {4, 7, 5}}
	matA := new(Matrix).Init(a)
	b := Data{{-0.06193087, -0.76241474, 0.64411826}, {0.9184855, -0.29608033, -0.26214659}, {0.39057517, 0.57537831, 0.71860339}}
	eig_vec := EigenVector33(matA, EigenValues33(matA))
	for i := range b {
		if !VEqual(&(b[i]), eig_vec.Row(i)) && !VEqual(&(b[i]), eig_vec.Row(i).MulNum(-1)) { // Discard sign difference
			fmt.Println(&(b[i]), eig_vec.Row(i))
			t.Fail()
		}
	}
}

func TestEigen33(t *testing.T) {
	a := Data{{1, 3, 4}, {3, 2, 7}, {4, 7, 5}}
	matA := new(Matrix).Init(a)
	b := Data{{-0.06193087, -0.76241474, 0.64411826}, {0.9184855, -0.29608033, -0.26214659}, {0.39057517, 0.57537831, 0.71860339}}
	eig_val, eig_vec := Eigen33(matA)
	if !VEqual(eig_val, &Vector{-3.67018839, -1.10871847, 12.77890686}) {
		t.Fail()
	}

	for i := range b {
		if !VEqual(&(b[i]), eig_vec.Row(i)) && !VEqual(&(b[i]), eig_vec.Row(i).MulNum(-1)) {
			t.Fail()
		}
	}
}

func BenchmarkEigenDecompose(b *testing.B) {
	b.Run("symmetric: size-3x3", func(b *testing.B) {
		m := GenerateRandomSymmetric33Matrix()
		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			EigenDecompose(m)
		}
	})

	for k := 1.0; k <= 2; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				EigenDecompose(m)
			}
		})
	}
}

func BenchmarkEigen33(b *testing.B) {
	b.Run("symmetric: size-3x3", func(b *testing.B) {
		m := GenerateRandomSymmetric33Matrix()
		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			Eigen33(m)
		}
	})
}
