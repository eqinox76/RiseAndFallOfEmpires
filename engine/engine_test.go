package engine

import (
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
	"testing"
	"time"
)

// let three times 1000 ships fight until no one is remaining or the planet gets occupied by a fleet
func TestFightsEnd(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()
	ge := GameEngine{Space: &space}

	passiveEmpire := space.CreateEmpire()
	passiveEmpire.Passive = true

	p := space.CreatePlanet(passiveEmpire)

	for eId := 0; eId < 3; eId++ {
		e := space.CreateEmpire()
		f := space.CreateFleet(p, e)
		f.LightSquads = 100
	}

	for len(p.Fleets) > 0 && p.Empire == passiveEmpire {
		ge.Step()
		for id, fleet := range p.Fleets{
			fmt.Println(id, fleet.Size())
		}
		fmt.Println("==", p.Empire, "==")
	}
}

func TestMoveAndFight(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()
	ge := GameEngine{Space: &space}

	e1:= space.CreateEmpire()
	e1.Passive = true
	p1 := space.CreatePlanet(e1)

	f1 := space.CreateFleet(p1, e1)
	f1.HeavySquads = 50

	e2:= space.CreateEmpire()
	p2 := space.CreatePlanet(e2)
	f2 := space.CreateFleet(p2, e2)
	f2.RangedSquads = 5

	for p1.Empire != p2.Empire {
		for _, fleet := range p1.Fleets {
			c := commands.MoveCommand{Destination: p2, Fleet:fleet}
			c.Execute()
		}

		ge.Step()

		for _, fleet := range p2.Fleets {
			fmt.Println("on p2", fleet.Size(), "ships", &fleet.Empire)
		}
		fmt.Println("==", &p1.Empire, "==", &p2.Empire, "==")

	}
}

func TestControl(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	ge := GameEngine{Space: &space}

	e1 := space.CreateEmpire()
	p := space.CreatePlanet(e1)
	p.Control = 0.5

	ge.Step()
	if p.Control <= 0.5 {
		t.Errorf("Control did not increase %v", p)
		return
	}

	space.CreateFleet(p, e1)

	old := p.Control
	ge.Step()
	if p.Control <= old {
		t.Errorf("Control did not increase %v", p)
		return
	}
}

func TestInvade(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()
	ge := GameEngine{Space: &space}

	e1 := space.CreateEmpire()
	e2 := space.CreateEmpire()
	p := space.CreatePlanet(e1)
	p.Control = 0.5

	f2:=space.CreateFleet(p, e2)
	f2.LightSquads = 10

	ge.Step()
	if p.Control >= 0.5 {
		t.Errorf("Control did not decrease %v", p)
		return
	}

	counter := 0
	for p.Empire == e1 {
		ge.Step()
		counter++
	}

	if p.Empire != e2 {
		t.Errorf("planet did not change owner %v", p)
		return
	}

	if counter < 10 {
		t.Errorf("planet changed owner to fast %v %d", p, counter)
		return
	}
}
