#

## Golina

A toy repo for reviewing linear algebra.

Realized functions up to now:

Matrix Operations: `Add`, `Sub`, `Mul`, `MulNum`, `Pow`, `Trace`, `T`, `Rank`, `Det`, `Adj`, `Inverse`, `Norm`, 
`Flat`, `GetSubMatrix`, `SetSubMatrix`
 
Eigen-Decomposition: `Eigen`, `EigenValues`, `EigenVector`

LU-Decomposition: `LUPDecompose`, `LUPSolve`, `LUPInvert`, `LUPDeterminant`, `LUPRank`

Vector `VEqual`, `Add`, `AddNum`, `Sub`, `SubNum`, `MulNum`, `Dot`, `Cross`, `SquareSum`, `Norm`, `Normalize`, 
`Convolve`

Benchmark:

```bash
CPU, 64-bit Linux
Intel: Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz, 32073 MB, Max threads(8)

BenchmarkLUPDeterminant/size-10-8                  20000             93133 ns/op
BenchmarkLUPDeterminant/size-100-8                  1000           1918208 ns/op
BenchmarkLUPDeterminant/size-1000-8                  100         621303080 ns/op
BenchmarkVector_SquareSum/size-10-8               200000              9216 ns/op
BenchmarkVector_SquareSum/size-100-8              100000             12798 ns/op
BenchmarkVector_SquareSum/size-1000-8              30000             58889 ns/op
BenchmarkVector_SquareSum/size-10000-8              3000            435305 ns/op
BenchmarkVector_SquareSum/size-100000-8              300           4876909 ns/op
BenchmarkEigen-8                                  100000             17866 ns/op
BenchmarkConvolve/size-10-8                       100000             20842 ns/op
BenchmarkConvolve/size-100-8                       30000             41146 ns/op
BenchmarkConvolve/size-1000-8                       3000            460712 ns/op
PASS
ok      golina  83.338s
```
