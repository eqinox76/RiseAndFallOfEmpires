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

	if g.GraphSize(g[0]) != 3 {
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

	if len(path) != 3 {
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []uint32{0, 4, 3}) {
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

	if len(path) != 4 {
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []uint32{3, 0, 1, 4}) {
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
	for i := uint32(0); i < 180; i++ {
		planets[i] = &pb.Planet{Id: i}
	}

	planets[173].Connected = append(planets[173].Connected, 60, 7)
	planets[7].Connected = append(planets[7].Connected, 178, 173, 50, 167)
	planets[178].Connected = append(planets[178].Connected, 36, 7, 41, 50)
	planets[50].Connected = append(planets[50].Connected, 7, 178)
	planets[167].Connected = append(planets[167].Connected, 39, 7)
	planets[36].Connected = append(planets[36].Connected, 178, 100)
	planets[39].Connected = append(planets[39].Connected, 156, 167)
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

	check(60, 36)
	check(178, 60)
	check(7, 36)
}

func TestShortestPathRealWorld2(t *testing.T) {
	//22 [107 83]
	//15 [59 129 162]
	//[15 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176 162 83 107 130 176]
	//	15 [59 129 162]
	//	162 [55 83 15 176]
	//	83 [162 107 22]
	//	107 [22 83 130]
	//	130 [107 122 176]
	planets := make(map[uint32]*pb.Planet)
	for i := uint32(0); i < 180; i++ {
		planets[i] = &pb.Planet{Id: i}
	}

	planets[22].Connected = append(planets[22].Connected, 107, 83)
	planets[83].Connected = append(planets[83].Connected, 162, 107, 22)
	planets[107].Connected = append(planets[107].Connected, 22, 83, 130)
	planets[130].Connected = append(planets[130].Connected, 107, 122, 176)
	planets[162].Connected = append(planets[162].Connected, 55, 83, 15, 176)
	planets[15].Connected = append(planets[15].Connected, 59, 129, 162)

	g := NewGraph(planets)

	path := g.ShortestPath(22, 15, false)

	if ! reflect.DeepEqual(path, []uint32{22, 83, 162, 15}) {
		t.Error(path, "expected to be [22, 83, 162, 15]")
	}
}

func TestShortestPathRealWorld3(t *testing.T) {
	//52 [95 73 42]
	//158 [29 199 17]
	//	[158 29 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123 73 123]
	//	158 [29 199 17]
	//	29 [158 123]
	//	123 [73 29]
	//	73 [123 52 95]
	//	123 [73 29]
	planets := make(map[uint32]*pb.Planet)
	for i := uint32(0); i < 180; i++ {
		planets[i] = &pb.Planet{Id: i}
	}

	planets[52].Connected = append(planets[52].Connected, 95, 73, 42)
	planets[123].Connected = append(planets[123].Connected, 73, 29)
	planets[73].Connected = append(planets[73].Connected, 123, 52, 95)
	planets[29].Connected = append(planets[29].Connected, 158, 123)
	planets[158].Connected = append(planets[158].Connected, 29, 199, 17)
	planets[95].Connected = append(planets[95].Connected, 52, 73)

	g := NewGraph(planets)

	check := func(start uint32, dest uint32) {
		path := g.ShortestPath(start, dest, false)
		//fmt.Println(start, dest, path)

		path = g.ShortestPath(52, 158, false)

		if ! reflect.DeepEqual(path, []uint32{52, 73, 123, 29, 158}) {
			t.Error(path, "expected to be [52 73 123 29 158]")
		}
	}

	check(95, 158)
	check(73, 158)
	check(123, 158)
}
