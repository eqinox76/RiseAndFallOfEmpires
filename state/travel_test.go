package state

import (
	"testing"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"fmt"
)

// test if the cycle detection works
func TestCycle(t *testing.T) {
	planets := make([] *pb.Planet, 3)
	planets[0] = &pb.Planet{Id: 0}
	planets[0].Connected = append(planets[0].Connected, 1, 2)
	planets[1] = &pb.Planet{Id: 1}
	planets[1].Connected = append(planets[1].Connected, 0)
	planets[2] = &pb.Planet{Id: 2}
	planets[2].Connected = append(planets[2].Connected, 0)

	g := NewGraph(planets)

	if g.HasCycle(planets[0]) || g.HasCycle(planets[1]) || g.HasCycle(planets[2]) {
		t.Error("Expected no Cycle ", g.HasCycle(planets[0]), g.HasCycle(planets[1]), g.HasCycle(planets[2]))
	}

	planets[1].Connected = append(planets[1].Connected, 2)
	planets[2].Connected = append(planets[2].Connected, 1)

	g = NewGraph(planets)

	if !g.HasCycle(planets[0]) || !g.HasCycle(planets[1]) ||!g.HasCycle(planets[2]) {
		fmt.Println(planets[0])
		fmt.Println(planets[1])
		fmt.Println(planets[2])
		t.Error("Expected Cycle ")
	}
}
