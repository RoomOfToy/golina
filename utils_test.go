package golina

import (
	"reflect"
	"testing"
)

func TestTernary(t *testing.T) {
	a, b := 1, 2
	if !(Ternary(a < b, a, b) == a) {
		t.Fail()
	}
}

func TestMEqual(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if !MEqual(matA, matA) {
		t.Fail()
	}
}

func TestVEqual(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{1, 2, 3}
	if !VEqual(v1, v2) {
		t.Fail()
	}
}

func TestMap(t *testing.T) {
	v := GenerateRandomVector(10)
	m := Map(VectorIter(v), func(item interface{}) interface{} {
		return item.(float64)
	})
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		for i := 0; i < s.Len(); i++ {
			if reflect.ValueOf(s.Index(i).Interface()).Convert(reflect.TypeOf(float64(0))).Float() != v.At(i) {
				t.Fail()
			}
		}
	}
}

func TestReduce(t *testing.T) {
	v := GenerateRandomVector(10)
	m := Reduce(VectorIter(v), func(item0, item1 interface{}) interface{} {
		return item0.(float64) + item1.(float64)
	})
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		if reflect.ValueOf(s.Interface()).Convert(reflect.TypeOf(float64(0))).Float() != v.Sum() {
			t.Fail()
		}
	}
}

func TestFilter(t *testing.T) {
	v := GenerateRandomVector(10)
	m := Filter(VectorIter(v), func(item interface{}) bool {
		return item.(float64) > 0
	})
	vv := make(Vector, 0, 10)
	for _, j := range *v {
		if j > 0 {
			vv = append(vv, j)
		}
	}
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		for i := 0; i < s.Len(); i++ {
			if reflect.ValueOf(s.Index(i).Interface()).Convert(reflect.TypeOf(float64(0))).Float() != vv.At(i) {
				t.Fail()
			}
		}
	}
}

func TestVectorIter(t *testing.T) {
	v := GenerateRandomVector(10)
	m := Map(VectorIter(v), func(item interface{}) interface{} {
		return item.(float64) * 2
	})
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		for i := 0; i < s.Len(); i++ {
			if reflect.ValueOf(s.Index(i).Interface()).Convert(reflect.TypeOf(float64(0))).Float() != (v.At(i) * 2) {
				t.Fail()
			}
		}
	}
}

func TestMatrixRowIter(t *testing.T) {
	mat := GenerateRandomMatrix(10, 3)
	m := Map(MatrixRowIter(mat), func(item interface{}) interface{} {
		// https://stackoverflow.com/questions/44543374/cannot-take-the-address-of-and-cannot-call-pointer-method-on
		v := item.(Vector)
		return v.Sum()
	})
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		for i := 0; i < s.Len(); i++ {
			v := mat.Row(i)
			if reflect.ValueOf(s.Index(i).Interface()).Convert(reflect.TypeOf(float64(0))).Float() != v.Sum() {
				t.Fail()
			}
		}
	}
}

func TestMatrixElementIter(t *testing.T) {
	mat := GenerateRandomMatrix(10, 3)
	m := Map(MatrixElementIter(mat), func(item interface{}) interface{} {
		return item.(float64) * 2
	})
	switch reflect.TypeOf(m).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(m)
		v := mat.Flat()
		for i := 0; i < s.Len(); i++ {
			if reflect.ValueOf(s.Index(i).Interface()).Convert(reflect.TypeOf(float64(0))).Float() != (v.At(i) * 2) {
				t.Fail()
			}
		}
	}
}

func TestLoad3DToMatrix(t *testing.T) {
	_, err := Load3DToMatrix("data.txt")
	if err != nil {
		println(err.Error())
		t.Fail()
	}
}

func TestWriteMatrixToTxt(t *testing.T) {
	mat := GenerateRandomMatrix(10, 10)
	err := WriteMatrixToTxt("test.txt", mat)
	if err != nil {
		println(err.Error())
		t.Fail()
	}
}
