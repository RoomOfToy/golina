package golina

import (
	"math"
)

func PointToPointDistance(p1, p2 *Vector) float64 {
	return p1.Sub(p2).Norm()
}

func PointToLineDistance(pt, linePt, lineDir *Vector) float64 {
	return pt.Sub(linePt).Sub(lineDir.MulNum(lineDir.Dot(pt.Sub(linePt)))).Norm()
}

func PointToPlaneDistance(pt, planeCenter, planeNormal *Vector) float64 {
	return math.Abs(planeNormal.Dot(pt.Sub(planeCenter)))
}

// https://en.wikipedia.org/wiki/Hausdorff_distance
type HausdorffDistance struct {
	distance       float64
	lIndex, rIndex int
}

func DirectedHausdorffDistance(pts1, pts2 *Matrix) *HausdorffDistance {
	r1, c1 := pts1.Dims()
	r2, c2 := pts2.Dims()
	if c1 != c2 {
		panic("points should have same coordinates")
	}
	cMax, cMin, d := 0., 0., 0.
	i, j, k := 0, 0, 0
	iStore, jStore, iRet, jRet := 0, 0, 0, 0
	noBreakOccurred := false
	hd := HausdorffDistance{0., 0, 0}
	for i = 0; i < r1; i++ {
		noBreakOccurred = true
		cMin = math.Inf(1)
		for j = 0; j < r2; j++ {
			d = 0.
			for k = 0; k < c1; k++ {
				d += (pts1.At(i, k) - pts2.At(j, k)) * (pts1.At(i, k) - pts2.At(j, k))
			}
			if d < cMax {
				noBreakOccurred = false
				break
			}
			if d < cMin {
				cMin = d
				iStore, jStore = i, j
			}
		}
		if cMin != math.Inf(1) && cMin > cMax && noBreakOccurred {
			cMax = cMin
			iRet = iStore
			jRet = jStore
		}
	}
	hd.distance = math.Sqrt(cMax)
	hd.lIndex = iRet
	hd.rIndex = jRet
	return &hd
}
