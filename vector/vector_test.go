package vector

import (
	"testing"
)

func TestAverage(t *testing.T) {
	a := Vec{10,10}
	b := Vec{11,10}

	if a.Dist(b) != 1 {
		t.Error("Expected 1, got ", a.Dist(b))
	}
}