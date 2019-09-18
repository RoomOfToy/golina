package pcl2mesh

import (
	"log"
	"sync"
)

// Singleton
type initVar struct {
	voxelSize                 []float64
	searchRadius              int // kNN
	normalSimilarityThreshold float64
	badCnt                    int
	planeMinPointNum          int
	planeMaxMSE               float64
}

var singleton *initVar
var once sync.Once

// Get init vars for point normal estimation
func GetInitVar() *initVar {
	once.Do(func() {
		singleton = &initVar{
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

func (v *initVar) GetVoxelSize() []float64 {
	return v.voxelSize
}

func (v *initVar) SetVoxelSize(voxelSize []float64) {
	if len(voxelSize) == 3 {
		v.voxelSize = voxelSize
	} else {
		log.Println("voxel size should have three dimensions")
	}
}

func (v *initVar) GetSearchRadius() int {
	return v.searchRadius
}

func (v *initVar) SetSearchRadius(radius int) {
	if radius > 0 {
		v.searchRadius = radius
	} else {
		log.Println("point normal calculation search radius should > 0")
	}
}

func (v *initVar) GetNormalSimilarityThreshold() float64 {
	return v.normalSimilarityThreshold
}

func (v *initVar) SetNormalSimilarityThreshold(similarityThreshold float64) {
	if similarityThreshold > 0 {
		v.normalSimilarityThreshold = similarityThreshold
	} else {
		log.Println("plain normal similarity threshold should > 0")
	}
}

func (v *initVar) GetBadCount() int {
	return v.badCnt
}

func (v *initVar) SetBadCount(badCnt int) {
	if badCnt > 0 {
		v.badCnt = badCnt
	} else {
		log.Println("plain bad condition count (edge case) should > 0")
	}
}

func (v *initVar) GetPlaneMinPointNum() int {
	return v.planeMinPointNum
}

func (v *initVar) SetPlaneMinPointNum(planeMinPointNum int) {
	if planeMinPointNum >= 3 {
		v.planeMinPointNum = planeMinPointNum
	} else {
		log.Println("plain minimum point number should >= 3")
	}
}

func (v *initVar) GetPlaneMaxMSE() float64 {
	return v.planeMaxMSE
}

func (v *initVar) SetPlaneMaxMSE(planeMaxMSE float64) {
	if planeMaxMSE >= 0 {
		v.planeMaxMSE = planeMaxMSE
	} else {
		log.Println("plain maximum MSE should >= 0")
	}
}
