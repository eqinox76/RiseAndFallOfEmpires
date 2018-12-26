package state

import (
	"fmt"
	"reflect"
	"testing"
)

func makePlanets(amount int) []*Planet {

	planets := make([]*Planet, amount)

	for i := range planets {
		var p Planet
		planets[i] = &p
	}

	return planets
}

// test if the cycle detection works
func TestCycle(t *testing.T) {
	planets := makePlanets(3)
	planets[0].Connected = append(planets[0].Connected, planets[1], planets[2])
	planets[1].Connected = append(planets[1].Connected, planets[0])
	planets[2].Connected = append(planets[2].Connected, planets[0])

	g := NewGraph(planets)

	if g.GraphSize(g.nodes[0]) != 3 {
		t.Error("Wrong size!", g.GraphSize(g.nodes[0]))
	}

	if g.HasCycle(planets[0]) || g.HasCycle(planets[1]) || g.HasCycle(planets[2]) {
		t.Error("Expected no Cycle ", g.HasCycle(planets[0]), g.HasCycle(planets[1]), g.HasCycle(planets[2]))
	}

	planets[1].Connected = append(planets[1].Connected, planets[2])
	planets[2].Connected = append(planets[2].Connected, planets[1])

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

	planets := makePlanets(5)

	planets[0].Connected = append(planets[0].Connected, planets[1], planets[4])
	planets[1].Connected = append(planets[1].Connected, planets[0], planets[2])
	planets[2].Connected = append(planets[2].Connected, planets[1], planets[3])
	planets[3].Connected = append(planets[2].Connected, planets[2], planets[4])
	planets[4].Connected = append(planets[2].Connected, planets[0], planets[3])

	g := NewGraph(planets)

	path := g.ShortestPath(planets[0], planets[3], false)

	if len(path) != 3 {
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []*Planet{planets[0], planets[4], planets[3]}) {
		t.Error(path, "expected to be [0 4 3]")
	}
}

func TestShortestPathWithLoop(t *testing.T) {
	// 0 - 1 - 4
	// |   |
	// 3 - 2

	planets := makePlanets(5)

	planets[0].Connected = append(planets[0].Connected, planets[1], planets[3])
	planets[1].Connected = append(planets[1].Connected, planets[0], planets[2], planets[4])
	planets[2].Connected = append(planets[2].Connected, planets[1], planets[3])
	planets[3].Connected = append(planets[3].Connected, planets[0], planets[2])
	planets[4].Connected = append(planets[3].Connected, planets[1])

	g := NewGraph(planets)

	path := g.ShortestPath(planets[3], planets[4], false)

	if len(path) != 4 {
		t.Error(path)
	}

	if ! reflect.DeepEqual(path, []*Planet{planets[3], planets[0], planets[1], planets[4]}) {
		t.Error(path, "expected to be [3 0 1 4]")
	}

}

func TestShortestPathRealWorld(t *testing.T) {
	//         50
	//        /  \
	// 173 - 7 - 178 - 36 - 100
	// |     |
	// 60   167 - 39 - 156

	planets := makePlanets(180)

	planets[173].Connected = append(planets[173].Connected, planets[60], planets[7])
	planets[7].Connected = append(planets[7].Connected, planets[178], planets[173], planets[50], planets[167])
	planets[178].Connected = append(planets[178].Connected, planets[36], planets[7], planets[41], planets[50])
	planets[50].Connected = append(planets[50].Connected, planets[7], planets[178])
	planets[167].Connected = append(planets[167].Connected, planets[39], planets[7])
	planets[36].Connected = append(planets[36].Connected, planets[178], planets[100])
	planets[39].Connected = append(planets[39].Connected, planets[156], planets[167])
	planets[60].Connected = append(planets[60].Connected, planets[173])

	g := NewGraph(planets)
	path := g.ShortestPath(planets[50], planets[39], false)

	if ! reflect.DeepEqual(path, []*Planet{planets[50], planets[7], planets[167], planets[39]}) {
		t.Error(path, "expected to be [50 7 167 39]")
	}

	stop := false

	go func() {
		innerGraph := NewGraph(planets)

		for !stop {
			path := innerGraph.ShortestPath(planets[36], planets[178], false)

			if ! reflect.DeepEqual(path, []*Planet{planets[36], planets[178]}) {
				t.Error(path, "expected to be [36 178]")
			}
		}
	}()

	check := func(start int, dest int) {
		g.ShortestPath(planets[start], planets[dest], false)
	}

	for i := 0; i < 10000; i++ {
		check(173, 36)
		//check(178, 60)
		//check(7, 36)
	}

	stop = true
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
	planets := makePlanets(200)

	planets[22].Connected = append(planets[22].Connected, planets[107], planets[83])
	planets[83].Connected = append(planets[83].Connected, planets[162], planets[107], planets[22])
	planets[107].Connected = append(planets[107].Connected, planets[22], planets[83], planets[130])
	planets[130].Connected = append(planets[130].Connected, planets[107], planets[122], planets[176])
	planets[162].Connected = append(planets[162].Connected, planets[55], planets[83], planets[15], planets[176])
	planets[15].Connected = append(planets[15].Connected, planets[59], planets[129], planets[162])

	g := NewGraph(planets)

	path := g.ShortestPath(planets[22], planets[15], false)

	if ! reflect.DeepEqual(path, []*Planet{planets[22], planets[83], planets[162], planets[15]}) {
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
	planets := makePlanets(200)

	planets[52].Connected = append(planets[52].Connected, planets[95], planets[73], planets[42])
	planets[123].Connected = append(planets[123].Connected, planets[73], planets[29])
	planets[73].Connected = append(planets[73].Connected, planets[123], planets[52], planets[95])
	planets[29].Connected = append(planets[29].Connected, planets[158], planets[123])
	planets[158].Connected = append(planets[158].Connected, planets[29], planets[199], planets[17])
	planets[95].Connected = append(planets[95].Connected, planets[52], planets[73])

	stop := false
	g := NewGraph(planets)

	go func() {
		g := NewGraph(planets)

		for !stop {
			path := g.ShortestPath(planets[52], planets[158], false)

			if ! reflect.DeepEqual(path, []*Planet{planets[52], planets[73], planets[123], planets[29], planets[158]}) {
				t.Error(path, "expected to be [52 73 123 29 158]")
			}
		}
	}()

	check := func(start int, dest int) {
		path := g.ShortestPath(planets[start], planets[ dest], false)
		//fmt.Println(start, dest, path)
		if len(path) <= 0 {
			t.Errorf("Path of length 0 found")
		}
	}

	for i := 0; i < 100000; i++ {
		check(95, 95)
		//check(95, 158)
		//check(73, 158)
		//check(123, 158)
	}

	stop = true

}
