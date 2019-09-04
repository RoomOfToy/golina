#

## Golina

A toy repo for reviewing linear algebra.

[doc](https://godoc.org/github.com/Harold2017/golina)

Pure Golang, No Dependencies.

Realized functions up to now:

- Matrix Operations: `Add`, `AddNum`, `Sub`, `Mul`, `MulVec`, `MulNum`, `Pow`, `Trace`, `T`, `Rank`, `Det`, `Adj`, `Inverse`, 
`Norm`, `Flat`, `GetSubMatrix`, `SetSubMatrix`, `SumCol`, `SumRow`, `Sum`, `Mean`, `CovMatrix`, `IsSymmetric`
- Eigen-Decomposition: `EigenDecompose`, `Eigen33`, `EigenValues33`, `EigenVector33`
- LU-Decomposition: `LUPDecompose`, `LUPSolve`, `LUPInvert`, `LUPDeterminant`, `LUPRank`
- QR-Decomposition: `Householder`, `QRDecomposition`
- Cholesky-Decomposition: `CholeskyDecomposition`
- SVD: `SVD`
- Matrix Transform: `Stretch`, `Rotate2D`, `Rotate3D`, `Translate`, `Shear2D`, `Shear3D`, 
`TransformOnRow` (for custom transform matrix), `ToAffineMatrix`, `Kabsch` (Superimpose)
- Vector Operations: `Add`, `AddNum`, `Sub`, `SubNum`, `MulNum`, `Dot`, `OuterProduct`, `Cross`, `SquareSum`, `Norm`, 
`Normalize`, `ToMatrix`, `Sum`, `AbsSum`, `Mean`, `Tile`, `Convolve`, `Max`, `Min`
- Distances: `PointToPointDistance`, `PointToLineDistance`, `PointToPlaneDistance`, `DirectedHausdorffDistance`; 
`TaxicabDistance`, `EuclideanDistance`, `SquaredEuclideanDistance`, `MinkowskiDistance`, `ChebyshevDistance`, 
`HammingDistance`, `CanberraDistance`
- k-Nearest-Neighbors: `KNearestNeighbor`, `KNearestNeighborsWithDistance` (work with above distance functions)
- k-Means: `KMeans`, `RandomMeans`, `KMeansPP`, `PPMeans`
- normal-Estimation: `PlanePCA`, `PlaneLinearSolveWeighted`
- Octree (concept only): `Octree` (based on hash map), `OctreeNode` (location encoded as map key)
- KD-Tree: `Insert`, `Search`, `FindMinValue`, `FindMinNode`, `DeleteNode`
- Utils Functions: `FloatEqual`, `MEqual`(matrix), `VEqual`(vector), `Ternary`, `String`(matrix, vector pretty-print), 
`Map`, `Reduce`, `Filter` (`Map`, `Reduce`, `Filter` here are just for tests, if you want to use it, you'd better change 
them from using `interface` with `reflect` module to `[]float64` for performance, since you have known the data type...), 
`Load3DToMatrix`, `WriteMatrixToTxt`

Benchmark(simple parallel `Mul`, need more optimization):

```bash
CPU, 64-bit Linux
Intel: Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz, 32073 MB, Max threads(8)

BenchmarkEigenDecompose/symmetric:_size-3x3-8    1000000              1267 ns/op
BenchmarkEigenDecompose/size-10-8                  50000             25288 ns/op
BenchmarkEigenDecompose/size-100-8                   100          14075608 ns/op
BenchmarkEigen33/symmetric:_size-3x3-8            300000              6117 ns/op
BenchmarkCholeskyDecomposition/size-10-4         1000000              1122 ns/op
BenchmarkCholeskyDecomposition/size-100-4           5000            340299 ns/op
BenchmarkCholeskyDecomposition/size-1000-4           100         320812408 ns/op
BenchmarkLUPDecompose/size-10-4                  1000000              1676 ns/op
BenchmarkLUPDecompose/size-100-4                    2000            712092 ns/op
BenchmarkLUPDecompose/size-1000-4                    100         770541089 ns/op
BenchmarkQRDecomposition/size-10-8                 10000            126480 ns/op
BenchmarkQRDecomposition/size-100-8                  100         460142966 ns/op
BenchmarkSVD/size-10-8                            100000             21107 ns/op
BenchmarkSVD/size-100-8                              200           9610030 ns/op
BenchmarkLUPDeterminant/size-10-8                1000000              1836 ns/op
BenchmarkLUPDeterminant/size-100-8                  2000            595163 ns/op
BenchmarkLUPDeterminant/size-1000-8                  100         556498134 ns/op
BenchmarkLUPRank/size-10-8                       1000000              1798 ns/op
BenchmarkLUPRank/size-100-8                         2000            596218 ns/op
BenchmarkLUPRank/size-1000-8                         100         557026660 ns/op
BenchmarkMatrix_Rank/size-10-8                    500000              2704 ns/op
BenchmarkMatrix_Rank/size-100-8                     1000           1720043 ns/op
BenchmarkMatrix_Rank/size-1000-8                     100        1754637844 ns/op
BenchmarkMatrix_Mul/size-10-8                     500000              3305 ns/op
BenchmarkMatrix_Mul/size-100-8                      1000           2167119 ns/op
BenchmarkMatrix_Mul/size-300-8                       100          72219939 ns/op
BenchmarkMatrix_Mul/size-400-8                       100         178311710 ns/op
BenchmarkMatrix_Mul/size-500-8                       100         366776279 ns/op
BenchmarkMatrix_MulNum/size-10-8                 2000000               935 ns/op
BenchmarkMatrix_MulNum/size-100-8                  30000             44814 ns/op
BenchmarkMatrix_MulNum/size-1000-8                   500           3011459 ns/op
BenchmarkVector_SquareSum/size-10-8            200000000              7.79 ns/op
BenchmarkVector_SquareSum/size-100-8            20000000              72.7 ns/op
BenchmarkVector_SquareSum/size-1000-8            2000000               762 ns/op
BenchmarkConvolve/size-10-8                      2000000               978 ns/op
BenchmarkConvolve/size-100-8                      100000             13134 ns/op
BenchmarkConvolve/size-1000-8                       5000            283011 ns/op
BenchmarkRotate3D/size-10x3-4                     300000              5303 ns/op
BenchmarkRotate3D/size-100x3-4                     30000             52046 ns/op
BenchmarkRotate3D/size-1000x3-4                     3000            523652 ns/op
BenchmarkKNearestNeighbors/size-10x3-4           1000000              1643 ns/op
BenchmarkKNearestNeighbors/size-100x3-4           100000             16588 ns/op
BenchmarkKNearestNeighbors/size-1000x3-4            5000            233153 ns/op
BenchmarkKMeans/size-10x3-4                        50000             21644 ns/op
BenchmarkKMeans/size-100x3-4                        2000            773231 ns/op
BenchmarkKMeans/size-1000x3-4                        100          47704397 ns/op
BenchmarkKMeansPP/size-10x3-4                      50000             27485 ns/op
BenchmarkKMeansPP/size-100x3-4                      2000            595368 ns/op
BenchmarkKMeansPP/size-1000x3-4                      100          51592136 ns/op
BenchmarkPlanePCA/size-10x3-8                     300000              5055 ns/op
BenchmarkPlanePCA/size-100x3-8                    100000             18980 ns/op
BenchmarkPlanePCA/size-1000x3-8                    10000            147370 ns/op
BenchmarkPlaneLinearSolveWeighted/size-10x3-8     300000              3989 ns/op
BenchmarkPlaneLinearSolveWeighted/size-100x3-8    100000             18511 ns/op
BenchmarkPlaneLinearSolveWeighted/size-1000x3-8    10000            148080 ns/op
BenchmarkKabsch/size-10x3-8                       100000             17576 ns/op
BenchmarkKabsch/size-100x3-8                       20000             73274 ns/op
BenchmarkKabsch/size-1000x3-8                       2000            580883 ns/op

BenchmarkGenerateRandomSparseMatrix/size-100000-8  50000             28520 ns/op
BenchmarkGenerateRandomSparseMatrix/size-1000000-8  5000            236498 ns/op
BenchmarkGenerateRandomSparseMatrix/size-10000000-8 1000           2176181 ns/op
BenchmarkGenerateRandomSparseMatrix/size-100000000-8 100          22174043 ns/op
BenchmarkSparseMatrix_MulVec/size-100000-8         10000            161415 ns/op
BenchmarkSparseMatrix_MulVec/size-1000000-8         1000           1289630 ns/op
BenchmarkSparseMatrix_MulVec/size-10000000-8         300           5338996 ns/op
BenchmarkSparseMatrix_MulVec/size-100000000-8        100          63072924 ns/op
```
