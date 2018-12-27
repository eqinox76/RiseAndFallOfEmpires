package strategies

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"testing"
)

func TestSendToFront(t *testing.T) {
	space := state.EmptySpace()
	e1 := space.CreateEmpire()
	e2 := space.CreateEmpire()

	p1 := space.CreatePlanet(e1)
	p2 := space.CreatePlanet(e1)
	p3 := space.CreatePlanet(e1)
	p4 := space.CreatePlanet(e2)

	p1.Connected = append(p1.Connected, p2)
	p1.Control = 1
	p2.Connected = append(p2.Connected, p3)
	p2.Control = 1
	p3.Connected = append(p3.Connected, p4)
	p3.Control = 1

	f1 := space.CreateFleet(p1, e1)
	f1.LightSquads = 1

	strat := Distributed{}
	strat.Init(e1)

	cmds := strat.Commands(&space)
	cmds[0].Execute()
	if len(cmds) != 1 || len(p2.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
	cmds = strat.Commands(&space)
	cmds[0].Execute()
	if len(cmds) != 1 || len(p3.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
	cmds = strat.Commands(&space)
	cmds[0].Execute()
	if len(cmds) != 1 || len(p4.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
}
