#

## Golina

A toy repo for reviewing linear algebra.

Realized functions up to now:

- Matrix Operations: `Add`, `Sub`, `Mul`, `MulNum`, `Pow`, `Trace`, `T`, `Rank`, `Det`, `Adj`, `Inverse`, `Norm`, 
`Flat`, `GetSubMatrix`, `SetSubMatrix`, `SumCol`, `SumRow`, `Sum`, `Mean`, `CovMatrix`
- Eigen-Decomposition: `Eigen`, `EigenValues`, `EigenVector`
- LU-Decomposition: `LUPDecompose`, `LUPSolve`, `LUPInvert`, `LUPDeterminant`, `LUPRank`
- QR-Decomposition: `Householder`, `QRDecomposition`
- Vector Operations: `Add`, `AddNum`, `Sub`, `SubNum`, `MulNum`, `Dot`, `Cross`, `SquareSum`, `Norm`, `Normalize`, 
`ToMatrix`, `Sum`, `Mean`, `Tile`, `Convolve`
- Helper Functions: `FloatEqual`, `Equal`(matrix), `VEqual`(vector), `Ternary`, `String`(matrix, vector pretty-print)

Benchmark(simple parallel `Mul`, need more optimization):

```bash
CPU, 64-bit Linux
Intel: Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz, 32073 MB, Max threads(8)

BenchmarkLUPDeterminant/size-10-8                1000000              1836 ns/op
BenchmarkLUPDeterminant/size-100-8                  2000            595163 ns/op
BenchmarkLUPDeterminant/size-1000-8                  100         556498134 ns/op
BenchmarkLUPRank/size-10-8                       1000000              1798 ns/op
BenchmarkLUPRank/size-100-8                         2000            596218 ns/op
BenchmarkLUPRank/size-1000-8                         100         557026660 ns/op
BenchmarkQRDecomposition/size-10-8                 10000            126480 ns/op
BenchmarkQRDecomposition/size-100-8                  100         460142966 ns/op
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
BenchmarkMatrix_Rank/size-10-8                    500000              2704 ns/op
BenchmarkMatrix_Rank/size-100-8                     1000           1720043 ns/op
BenchmarkMatrix_Rank/size-1000-8                     100        1754637844 ns/op
BenchmarkEigen/size-3-8                           100000             15462 ns/op
BenchmarkConvolve/size-10-8                      2000000               978 ns/op
BenchmarkConvolve/size-100-8                      100000             13134 ns/op
BenchmarkConvolve/size-1000-8                       5000            283011 ns/op
```
