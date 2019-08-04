package strategies

import (
	"bytes"
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"testing"
)

// p1 - p2 - p3 - p4
// e1   e1   e1   e2
// a new fleet on p1 should move to p4
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
	cmds[0].Execute(&space)
	if len(cmds) != 1 || len(p2.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
	cmds = strat.Commands(&space)
	cmds[0].Execute(&space)
	if len(cmds) != 1 || len(p3.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
	cmds = strat.Commands(&space)
	cmds[0].Execute(&space)
	if len(cmds) != 1 || len(p4.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
}

// p1 - p2
// having two fleets on p1 should trigger one to move to p2
func TestOccupyOnlyWithOneFleet(t *testing.T) {
	space := state.EmptySpace()
	e1 := space.CreateEmpire()
	e2 := space.CreateEmpire()

	p1 := space.CreatePlanet(e1)
	p2 := space.CreatePlanet(e2)

	p1.Connected = append(p1.Connected, p2)
	p1.Control = 0.1

	f1 := space.CreateFleet(p1, e1)
	f1.LightSquads = 1
	f2 := space.CreateFleet(p1, e1)
	f2.LightSquads = 1

	strat := Distributed{}
	strat.Init(e1)

	cmds := strat.Commands(&space)
	cmds[0].Execute(&space)
	if len(cmds) != 1 || len(p1.Fleets) != 1 || len(p2.Fleets) != 1 {
		t.Error("Fleet did not move to next planet")
		return
	}
}

// p0 - p1 - p3
//  |   |
//    p2
func TestMovement(t *testing.T) {
	space := state.EmptySpace()
	e1 := space.CreateEmpire()
	e2 := space.CreateEmpire()

	p0 := space.CreatePlanet(e1)
	p1 := space.CreatePlanet(e2)
	p2 := space.CreatePlanet(e1)
	p3 := space.CreatePlanet(e2)

	p0.Connected = append(p0.Connected, p1)
	p0.Connected = append(p0.Connected, p2)
	p0.Control = 1
	p1.Connected = append(p1.Connected, p0)
	p1.Connected = append(p1.Connected, p2)
	p1.Connected = append(p1.Connected, p3)
	p1.Control = 1
	p2.Connected = append(p2.Connected, p0)
	p2.Connected = append(p2.Connected, p1)
	p2.Control = 1
	p2.Connected = append(p3.Connected, p1)
	p3.Control = 1

	fleet := space.CreateFleet(p0, e1)
	fleet.RangedSquads = 100
	fleet = space.CreateFleet(p2, e1)
	fleet.RangedSquads = 10
	fleet = space.CreateFleet(p2, e1)
	fleet.RangedSquads = 1

	strat := Distributed{}
	strat.Init(e1)

	log := bytes.Buffer{}

	for i:= 0; i < 10 ; i++ {
		for i, planet := range space.Planets{
			fleets := 0
			ships := 0
			for _, fleet := range planet.Fleets{
				if fleet.Empire == e1 {
					fleets++
					ships += fleet.Size()
				}
			}
			_, _ = fmt.Fprintln(&log, i, "\t", len(planet.Fleets), "\t", fleets, "\t", ships, "\t", planet.Control)
		}
		_, _ = fmt.Fprintln(&log,"id \t fleets \t e fleets \t e ships \t control")

		for _, cmd := range strat.Commands(&space) {
			cmd.Execute(&space)
		}
	}

	if len(space.Fleets) != 2{
		t.Errorf("Fleets not merged\n%s", log.String())
	}

	if len(p0.Fleets) != 0{
		t.Errorf("Fleets still on own planet\n%s", log.String())
	}

	if len(p2.Fleets) != 0{
		t.Errorf("Fleets still on own planet\n%s", log.String())
	}
}
