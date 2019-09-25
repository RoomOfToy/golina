package mesh

import (
	"log"
	"sync"
)

// Singleton
type initVarNormalEst struct {
	voxelSize                 []float64
	searchRadius              int // kNN
	normalSimilarityThreshold float64
	badCnt                    int
	planeMinPointNum          int
	planeMaxMSE               float64
}

var singleton *initVarNormalEst
var once sync.Once

// Get init vars for point normal estimation
func GetInitVarNormalEst() *initVarNormalEst {
	once.Do(func() {
		singleton = &initVarNormalEst{
			voxelSize:                 []float64{15, 15, 15},
			searchRadius:              5,
			normalSimilarityThreshold: 25,
			badCnt:                    3,
			planeMinPointNum:          3,
			planeMaxMSE:               10,
		}
	})
	return singleton
}

func (v *initVarNormalEst) GetVoxelSize() []float64 {
	return v.voxelSize
}

func (v *initVarNormalEst) SetVoxelSize(voxelSize []float64) {
	if len(voxelSize) == 3 {
		v.voxelSize = voxelSize
	} else {
		log.Println("voxel size should have three dimensions")
	}
}

func (v *initVarNormalEst) GetSearchRadius() int {
	return v.searchRadius
}

func (v *initVarNormalEst) SetSearchRadius(radius int) {
	if radius > 0 {
		v.searchRadius = radius
	} else {
		log.Println("point normal calculation search radius should > 0")
	}
}

func (v *initVarNormalEst) GetNormalSimilarityThreshold() float64 {
	return v.normalSimilarityThreshold
}

func (v *initVarNormalEst) SetNormalSimilarityThreshold(similarityThreshold float64) {
	if similarityThreshold > 0 {
		v.normalSimilarityThreshold = similarityThreshold
	} else {
		log.Println("plain normal similarity threshold should > 0")
	}
}

func (v *initVarNormalEst) GetBadCount() int {
	return v.badCnt
}

func (v *initVarNormalEst) SetBadCount(badCnt int) {
	if badCnt > 0 {
		v.badCnt = badCnt
	} else {
		log.Println("plain bad condition count (edge case) should > 0")
	}
}

func (v *initVarNormalEst) GetPlaneMinPointNum() int {
	return v.planeMinPointNum
}

func (v *initVarNormalEst) SetPlaneMinPointNum(planeMinPointNum int) {
	if planeMinPointNum >= 3 {
		v.planeMinPointNum = planeMinPointNum
	} else {
		log.Println("plain minimum point number should >= 3")
	}
}

func (v *initVarNormalEst) GetPlaneMaxMSE() float64 {
	return v.planeMaxMSE
}

func (v *initVarNormalEst) SetPlaneMaxMSE(planeMaxMSE float64) {
	if planeMaxMSE >= 0 {
		v.planeMaxMSE = planeMaxMSE
	} else {
		log.Println("plain maximum MSE should >= 0")
	}
}

type initVarTriangulation struct {
	initHullVerticesCnt int
	initHullFacesCnt    int
	unitSphereRadius    float64
}

func GetInitVarTriangulation() *initVarTriangulation {
	return &initVarTriangulation{
		initHullVerticesCnt: 6,
		initHullFacesCnt:    8,
		unitSphereRadius:    1,
	}
}
