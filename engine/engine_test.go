package engine

import (
	"testing"
	"math/rand"
	"time"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)


// let three times 1000 ships figth until only one fleet is remaining
func TestFightsEnd(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p := space.CreateNewPlanet(space.CreateEmpire("doesntmatter"))

	for e_id := 0; e_id < 3; e_id++ {
		e := space.CreateEmpire("none")

		for i := 0; i < 1000; i++ {
			space.CreateShip(p, e)
		}
	}

	fleets := getFleets(space.Ships, p)
	for ; len(fleets) > 1; fleets = getFleets(space.Ships, p){
		//for id, fleet := range fleets{
		//	fmt.Println(id, len(fleet))
		//}
		//fmt.Println("=====")
		computeFight(&space, p)
	}
}
