package pcl2mesh

import (
	"log"
	"sync"
)

// Singleton
type initVar struct {
	voxelSize                 []float64
	searchRadius              float64
	normalSimilarityThreshold float64
	badCnt                    int
	planeMinPointNum          int
}

var singleton *initVar
var once sync.Once

// Get init vars for point normal estimation
func GetInitVar() *initVar {
	once.Do(func() {
		singleton = &initVar{
			voxelSize:                 []float64{5, 5, 5},
			searchRadius:              10,
			normalSimilarityThreshold: 25,
			badCnt:                    3,
			planeMinPointNum:          3,
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

func (v *initVar) GetSearchRadius() float64 {
	return v.searchRadius
}

func (v *initVar) SetSearchRadius(radius float64) {
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
