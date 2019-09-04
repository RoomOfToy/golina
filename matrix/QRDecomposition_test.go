package matrix

import (
	"math"
	"strconv"
	"testing"
)

func TestSign(t *testing.T) {
	if Sign(-1.) != -1 || Sign(1.) != 1 {
		t.Fail()
	}
}

func TestUnitVector(t *testing.T) {
	if !VEqual(UnitVector(3, 0).Col(0), &Vector{1, 0, 0}) {
		t.Fail()
	}
}

func TestQR(t *testing.T) {
	a := Data{{5, 43, 65}, {-76, 32, 12}, {4, -3, 2}}
	matA := new(Matrix).Init(a)
	q, r := QRDecomposition(matA)
	if !MEqual(q, new(Matrix).Init(Data{
		{-0.06555721150307149, -0.997423940968518, 0.02911587200775649},
		{0.9964696148466867, -0.06390513009925516, 0.05444668065450027},
		{-0.05244576920245719, 0.032582454324799455, 0.9980920924258115}})) ||
		!MEqual(r, new(Matrix).Init(Data{
			{-76.26925986267338, 29.22540488806927, 7.5915250920556785},
			{0, -45.03194098779684, -65.53425281549514},
			{0, 0, 4.5420760332096455}})) {
		t.Fail()
	}
}

func BenchmarkQRDecomposition(b *testing.B) {
	for k := 1.0; k <= 2; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				QRDecomposition(m)
			}
		})
	}
}
