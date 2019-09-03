package golina

import (
	"fmt"
	"testing"
)

func TestSparseMatrix_GetSubSparseMatrix(t *testing.T) {
	SA := GenerateRandomSparseMatrix(100, 100, 20)
	SubSA := SA.GetSubSparseMatrix(50, 50, 50, 50)
	fmt.Println(SubSA.Offset)
	fmt.Println(SA.FindFirstNonZeroInSubMatrix(SA.RowColToIndex(50, 50)))
}
