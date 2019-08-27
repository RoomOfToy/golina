package golina

// calculate plane norm following the `Total Least Squares`
// https://en.wikipedia.org/wiki/Total_least_squares

// common solution: `Principle Component Analysis`
// https://en.wikipedia.org/wiki/Principal_component_analysis
func PlanePCA(points *Matrix) *Vector {
	row, col := points.Dims()
	if col > 3 {
		panic("Only 3D points is supported")
	}
	if row < 3 {
		panic("Not enough points to fit a plane")
	}
	cov := points.CovMatrix()
	_, eigVec := Eigen(cov)
	return eigVec.Row(2)
}

// https://www.ilikebigbits.com/2017_09_25_plane_from_points_2.html
// ax + by +z + d = 0
// ax + by + d = -z
// [X, Y, 1] * [a, b, d].T() = -Z
// [X, Y, 1].T() * [X, Y, 1] * [a, b, d].T() = [X, Y, 1].T() * -Z
//
// [ Σxx, Σxy, Σx,        [ a,          [ Σxz,
//   Σyx, Σyy, Σyz,    *    b,    =   -   Σyz,
//   Σx,  Σy,  N   ]        d ]           0   ]
// define x, y, z is relative to the centroid (mean) of points, then Σx = Σy = Σz = 0, then N * d = 0, then d = 0
// [ Σxx, Σxy,   *  [ a,    =  - [ Σxz,
//   Σyx, Σyy ]       b ]          Σyz]
// according to Cramer's rule (https://en.wikipedia.org/wiki/Cramer%27s_rule)
// D = Σxx * Σyy - Σxy * Σxy
// a = (Σyz * Σxy - Σxz * Σyy) / D
// b = (Σxy * Σxz - Σxx * Σyz) / D
// so n = [a, b, 1].T(), multiply by D, then n = [a, b, D].T(), then normalize it we will get the plane norm
// but this works only when z-component of the plane normal is non-zero, if it is, then we can use the x or y component for calculation
// sine the above only minimize the squares of the residuals as perpendicular to the main axis, not the residuals perpendicular to the plane,
// then use the below weighted way to calculate the components of plain normal is more reasonable.
func PlaneLinearSolveWeighted(points *Matrix) *Vector {
	row, col := points.Dims()
	if col > 3 {
		panic("Only 3D points is supported")
	}
	if row < 3 {
		panic("Not enough points to fit a plane")
	}
	cov := points.CovMatrix()
	xx, xy, xz, yy, yz, zz := cov.At(0, 0), cov.At(0, 1), cov.At(0, 2), cov.At(1, 1), cov.At(1, 2), cov.At(2, 2)
	/*
		// calculate cov
		mean := points.Mean(0)
		for _, p := range points._array {
			r := p.Sub(mean)
			xx += r.At(0) * r.At(0)
			xy += r.At(0) * r.At(1)
			xz += r.At(0) * r.At(2)
			yy += r.At(1) * r.At(1)
			yz += r.At(1) * r.At(2)
			zz += r.At(2) * r.At(2)
		}
		xx = xx / float64(row)
		xy = xy / float64(row)
		xz = xz / float64(row)
		yy = yy / float64(row)
		yz = yz / float64(row)
		zz = zz / float64(row)

		fmt.Println(new(Matrix).Init(Data{{xx, xy, xz}, {xy, yy, yz}, {xz, yz, zz}}))
	*/

	weightedDir := Vector{0., 0., 0.}

	// x direction
	detX := yy*zz - yz*yz
	axis_dir := Vector{detX, xz*yz - xy*zz, xy*yz - xz*yy}
	weight := detX * detX
	if weightedDir.Dot(&axis_dir) < 0. {
		weight = -weight
	}
	weightedDir = *(weightedDir.Add(axis_dir.MulNum(weight)))

	// y direction
	detY := xx*zz - xz*xz
	axis_dir = Vector{xz*yz - xy*zz, detY, xy*xz - yz*xx}
	weight = detY * detY
	if weightedDir.Dot(&axis_dir) < 0. {
		weight = -weight
	}
	weightedDir = *(weightedDir.Add(axis_dir.MulNum(weight)))

	// z direction
	detZ := xx*yy - xy*xy
	axis_dir = Vector{xy*yz - xz*yy, xy*xz - yz*xx, detZ}
	weight = detZ * detZ
	if weightedDir.Dot(&axis_dir) < 0. {
		weight = -weight
	}
	weightedDir = *(weightedDir.Add(axis_dir.MulNum(weight)))

	planeNorm := weightedDir.Normalize()
	return planeNorm
}
