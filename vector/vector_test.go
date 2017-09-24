package vector

import (
	"testing"
	"math"
	"fmt"
)

func TestAverage(t *testing.T) {
	a := Vec{10,10}
	b := Vec{11,10}

	if a.Dist(b) != 1 {
		t.Error("Expected 1, got ", a.Dist(b))
	}
}

func TestDegree(t *testing.T){
	fmt.Println(math.Cos(90 ))
}