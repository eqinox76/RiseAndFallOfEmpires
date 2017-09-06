package engine

import (
	"testing"
	"math/rand"
	"time"
)

func TestAverage(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))
	var a uint32 = 1000
	var b uint32 = 1000

	// make sure fights end
	for a != 0 && b != 0{
		a, b = Fight(a,b)
	}
}