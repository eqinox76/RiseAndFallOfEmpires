package state

import (
	"testing"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"fmt"
	"reflect"
)

// test if the cycle detection works
func TestCycle(t *testing.T) {
	planets := make(map[uint32]*pb.Planet)
	planets[0] = &pb.Planet{Id: 0}
	planets[0].Connected = append(planets[0].Connected, 1, 2)
	planets[1] = &pb.Planet{Id: 1}
	planets[1].Connected = append(planets[1].Connected, 0)
	planets[2] = &pb.Planet{Id: 2}
	planets[2].Connected = append(planets[2].Connected, 0)

	g := NewGraph(planets)

	if g.GraphSize(g[0]) != 3{
		t.Error("Wrong size!", g.GraphSize(g[0]))
	}

	if g.HasCycle(planets[0]) || g.HasCycle(planets[1]) || g.HasCycle(planets[2]) {
		t.Error("Expected no Cycle ", g.HasCycle(planets[0]), g.HasCycle(planets[1]), g.HasCycle(planets[2]))
	}

	planets[1].Connected = append(planets[1].Connected, 2)
	planets[2].Connected = append(planets[2].Connected, 1)

	g = NewGraph(planets)

	if !g.HasCycle(planets[0]) || !g.HasCycle(planets[1]) || !g.HasCycle(planets[2]) {
		fmt.Println(planets[0])
		fmt.Println(planets[1])
		fmt.Println(planets[2])
		t.Error("Expected Cycle ")
	}
}

func TestShortestPath(t *testing.T) {
	//0 -> 1 -> 2 -> 3
	//  `->   4   ->'

	planets := make(map[uint32]*pb.Planet)
	planets[0] = &pb.Planet{Id: 0}
	planets[1] = &pb.Planet{Id: 1}
	planets[2] = &pb.Planet{Id: 2}
	planets[3] = &pb.Planet{Id: 3}
	planets[4] = &pb.Planet{Id: 4}

	planets[0].Connected = append(planets[0].Connected, 1, 4)
	planets[1].Connected = append(planets[1].Connected, 0, 2)
	planets[2].Connected = append(planets[2].Connected, 1, 3)
	planets[3].Connected = append(planets[2].Connected, 2, 4)
	planets[4].Connected = append(planets[2].Connected, 0, 3)

	g := NewGraph(planets)

	path := g.ShortestPath(0, 3, false)

	if len(path) != 3{
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []uint32 {0, 4, 3}){
		t.Error(path, "expected to be [0, 4, 3]")
	}
}


func TestShortestPathWithLoop(t *testing.T) {
	// 0 - 1 - 4
	// |   |
	// 3 - 2

	planets := make(map[uint32]*pb.Planet)
	planets[0] = &pb.Planet{Id: 0}
	planets[1] = &pb.Planet{Id: 1}
	planets[2] = &pb.Planet{Id: 2}
	planets[3] = &pb.Planet{Id: 3}
	planets[4] = &pb.Planet{Id: 4}

	planets[0].Connected = append(planets[0].Connected, 1, 3)
	planets[1].Connected = append(planets[1].Connected, 0, 2, 4)
	planets[2].Connected = append(planets[2].Connected, 1, 3)
	planets[3].Connected = append(planets[3].Connected, 0, 2)
	planets[4].Connected = append(planets[3].Connected, 1)

	g := NewGraph(planets)

	path := g.ShortestPath(3, 4, false)

	if len(path) != 4{
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []uint32 {3, 0, 1 ,4}){
		t.Error(path, "expected to be [3, 0, 1 ,4]")
	}

}

func TestShortestPathRealWorld(t *testing.T) {
	//         50
	//        /  \
	// 173 - 7 - 178 - 36
	// |     |
	// 60   167 - 39

	planets := make(map[uint32]*pb.Planet)
	planets[173] = &pb.Planet{Id: 173}
	for i:= uint32(0); i < 180; i++ {
		planets[i] = &pb.Planet{Id: i}
	}

	planets[173].Connected = append(planets[173].Connected, 60, 7)
	planets[7].Connected = append(planets[7].Connected, 178, 173, 50, 167)
	planets[178].Connected = append(planets[178].Connected, 36, 7, 41, 50)
	planets[50].Connected = append(planets[50].Connected, 7, 178)
	planets[167].Connected = append(planets[167].Connected, 39, 7)
	planets[36].Connected = append(planets[36].Connected, 178, 100)
	planets[39].Connected = append(planets[39].Connected, 156,167)
	planets[60].Connected = append(planets[60].Connected, 173)

	g := NewGraph(planets)

	check := func(start uint32, dest uint32) {
		path := g.ShortestPath(start, dest, false)
		//fmt.Println(start, dest, path)

		path = g.ShortestPath(50, 39, false)

		if ! reflect.DeepEqual(path, []uint32{50, 7, 167, 39}) {
			t.Error(path, "expected to be [50, 7, 167, 39]")
		}
	}

	check(60,36)
	check(178,60)
	check(7,36)

}
