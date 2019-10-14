package container

import (
	"math/rand"
	"time"
)

func GenerateRandomInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() - rand.Int()
}
