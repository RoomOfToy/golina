package stats

import "golina/matrix"

// Principal Component Analysis
//	calculate principal components direction vectors and corresponding column variances of scores
//	https://en.wikipedia.org/wiki/Principal_component_analysis
//	https://towardsdatascience.com/pca-and-svd-explained-with-numpy-5d13b0d2a4d8
//	SVD or Eigen (https://stats.stackexchange.com/questions/314046/why-does-andrew-ng-prefer-to-use-svd-and-not-eig-of-covariance-matrix-to-do-pca)
//	Here i prefer to use Eigen, since my Eigen implementation is faster than SVD due to covariance matrix is always real symmetric matrix
//	SVD way can be found in `PlanePcaSVD` of `spatial/normalEstimation.go`
func PrincipalComponents(dataSet *matrix.Matrix, weights *matrix.Vector) (pcs *matrix.Matrix, colVars *matrix.Vector) {
	row, col := dataSet.Dims()
	if weights != nil && weights.Length() != row {
		panic("length of weights vector should be equal to data matrix's rows")
	}
	// From wiki: https://en.wikipedia.org/wiki/Mahalanobis_distance
	// the ellipsoid that best represents the set's probability distribution can be estimated by building the covariance matrix of the samples
	cov := new(matrix.Matrix)
	if weights == nil {
		cov = dataSet.CovMatrix()
	} else {
		data := make(matrix.Data, row)
		for i, v := range dataSet.Data {
			data[i] = *(v.MulNum(weights.At(i)))
		}
		cov = new(matrix.Matrix).Init(data).CovMatrix()
	}
	eigVec, eigVal := matrix.EigenDecompose(cov) // eigVec is in ascending way
	tmpV := make(matrix.Vector, col)
	tmpM := matrix.ZeroMatrix(col, col)
	for i, v := range eigVec.T().Data {
		tmpM.Data[col-1-i] = v
		tmpV[i] = eigVal.At(col-1-i, col-1-i)
	}
	colVars = &tmpV
	pcs = tmpM.T()
	return
}
