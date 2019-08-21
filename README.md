#

## Golina

A toy repo for reviewing linear algebra.

Realized functions up to now:

Matrix Operations: `Add`, `Sub`, `Mul`, `MulNum`, `Pow`, `Trace`, `T`, `Rank`, `Det`, `Adj`, `Inverse`, `Norm`, 
`Flat`, `GetSubMatrix`, `SetSubMatrix`
 
Eigen-Decomposition: `Eigen`, `EigenValues`, `EigenVector`

LU-Decomposition: `LUPDecompose`, `LUPSolve`, `LUPInvert`, `LUPDeterminant`, `LUPRank`

QR-Decomposition: `Householder`, `QRDecomposition`

Vector `VEqual`, `Add`, `AddNum`, `Sub`, `SubNum`, `MulNum`, `Dot`, `Cross`, `SquareSum`, `Norm`, `Normalize`, 
`ToMatrix`, `Convolve`

Benchmark(No parallel code now, need optimize):

```bash
CPU, 64-bit Linux
Intel: Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz, 32073 MB, Max threads(8)
Used Core: 1, Single Thread

BenchmarkLUPDeterminant/size-10-8                  20000             93133 ns/op
BenchmarkLUPDeterminant/size-100-8                  1000           1918208 ns/op
BenchmarkLUPDeterminant/size-1000-8                  100         621303080 ns/op
BenchmarkMatrix_Det/size-10-8                        100        2121988644 ns/op
BenchmarkLUPRank/size-10-8                         20000             90408 ns/op
BenchmarkLUPRank/size-100-8                         1000           1988381 ns/op
BenchmarkLUPRank/size-1000-8                         100         663358579 ns/op
BenchmarkMatrix_Rank/size-10-8                     20000             91392 ns/op
BenchmarkMatrix_Rank/size-100-8                      500           2834004 ns/op
BenchmarkMatrix_Rank/size-1000-8                     100        1811608687 ns/op
BenchmarkEigen-8                                  100000             17866 ns/op
BenchmarkConvolve/size-10-8                       100000             20842 ns/op
BenchmarkConvolve/size-100-8                       30000             41146 ns/op
BenchmarkConvolve/size-1000-8                       3000            460712 ns/op
BenchmarkQRDecomposition/size-10-8                  5000            264976 ns/op
BenchmarkQRDecomposition/size-100-8                  100         610476336 ns/op
PASS
ok      golina  83.338s
```
