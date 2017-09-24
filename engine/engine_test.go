package engine

import (
	"testing"
	"math/rand"
	"time"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"fmt"
)

// let three times 1000 ships figth until no one is remaining or the planet gets occupied by a fleet
func TestFightsEnd(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p := space.CreatePlanet(space.CreateEmpire())

	for e_id := 0; e_id < 3; e_id++ {
		e := space.CreateEmpire()
		for i := 0; i < 1000; i++ {
			space.CreateShip(p, e)
		}
	}

	fleets := state.GetFleets(space.Ships, p)
	for len(p.Orbiting) > 0 && p.Empire == 0 {
		Step(&space)
		fleets = state.GetFleets(space.Ships, p)
		for id, fleet := range fleets {
			fmt.Println(id, len(fleet))
		}
		fmt.Println("==", p.Empire, "==")
	}
}

func TestMove(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p1 := space.CreatePlanet(space.CreateEmpire())
	for i := 0; i < 50; i ++ {
		space.CreateShip(p1, space.Empires[p1.Empire])
	}
	p2 := space.CreatePlanet(space.CreateEmpire())
	s := space.CreateShip(p2, space.Empires[p2.Empire])

	space.MoveShip(s.Id, p2.Id, p1.Id)

	if len(p2.Orbiting) > 0 {
		t.Errorf("%s ships still on p2", len(p2.Orbiting))
	}

	if len(p1.Orbiting) != 51 {
		t.Errorf("%d ships on p1", len(p1.Orbiting))
	}

	switch x := s.GetPosition().(type) {
	case *pb.Ship_Orbiting:
		planet := space.Planets[x.Orbiting]

		if planet.Id != p1.Id {
			t.Error("Wrong position", s)
		}
	default:
		t.Error(s, "Orbits no planet")
	}

	if len(state.GetFleets(space.Ships, p1)) != 2 {
		t.Error(state.GetFleets(space.Ships, p1))
	}
}

func TestDestroy(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p := space.CreatePlanet(space.CreateEmpire())
	for i := 0; i < 50; i ++ {
		space.CreateShip(p, space.Empires[p.Empire])
	}

	space.RemoveShip(space.Ships[rand.Uint64()%uint64(len(space.Ships)-1)])
	space.CreateShip(p, space.Empires[p.Empire])

	if len(p.Orbiting) != 50 {
		t.Errorf("%d ships on p1", len(p.Orbiting))
	}
}

func TestCreateAndDestroy(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p := space.CreatePlanet(space.CreateEmpire())

	s1 := space.CreateShip(p, space.Empires[p.Empire])
	s2 := space.CreateShip(p, space.Empires[p.Empire])
	s3 := space.CreateShip(p, space.Empires[p.Empire])
	s4 := space.CreateShip(p, space.Empires[p.Empire])

	space.RemoveShip(s2)
	space.RemoveShip(s3)
	s5 := space.CreateShip(p, space.Empires[p.Empire])
	for k, _ := range space.Ships {
		fmt.Print(k, "|")
	}

	if s1.Id == s5.Id ||
		s4.Id == s5.Id {
		t.Error(s1, s2, s3, s4, s5)
	}

	s2 = space.CreateShip(p, space.Empires[p.Empire])
	for k, _ := range space.Ships {
		fmt.Print(k, "|")
	}
	if s1.Id == s2.Id ||
		s4.Id == s2.Id ||
		s5.Id == s2.Id {
		t.Error(s1, s2, s3, s4, s5)
	}
	s3 = space.CreateShip(p, space.Empires[p.Empire])
	for k, _ := range space.Ships {
		fmt.Print(k, "|")
	}
	if s1.Id == s3.Id ||
		s2.Id == s3.Id ||
		s4.Id == s3.Id ||
		s5.Id == s3.Id {
		t.Error(s1, s2, s3, s4, s5)
	}
}

func TestMoveAndFight(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	p1 := space.CreatePlanet(space.CreateEmpire())
	space.Empires[p1.Empire].Passive = true
	for i := 0; i < 50; i ++ {
		space.CreateShip(p1, space.Empires[p1.Empire])
	}
	p2 := space.CreatePlanet(space.CreateEmpire())

	for p1.Empire != p2.Empire {
		for ship, _ := range p2.Orbiting {
			space.MoveShip(ship, p2.Id, p1.Id)
		}

		if len(p2.Orbiting) > 0 {
			t.Error(p2.Orbiting)
		}

		Step(&space)

		fleets := state.GetFleets(space.Ships, p1)
		for id, fleet := range fleets {
			fmt.Println(id, len(fleet))
		}
		fmt.Println("==", p1.Empire, "==")
	}
}

func TestControl(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	e1 := space.CreateEmpire()
	p := space.CreatePlanet(e1)
	p.Control = 0.5

	Step(&space)
	if p.Control <= 0.5 {
		t.Errorf("Control did not increase %s", p)
	}

	for i := 0; i < 10; i ++ {
		space.CreateShip(p, e1)
	}

	old := p.Control
	Step(&space)
	if p.Control <= old {
		t.Errorf("Control did not increase %s", p)
	}
}

func TestInvade(t *testing.T) {
	rand.Seed(int64(time.Now().Second()))

	space := state.EmptySpace()

	e1 := space.CreateEmpire()
	e2 := space.CreateEmpire()
	e3 := space.CreateEmpire()
	p := space.CreatePlanet(e1)
	p.Control = 0.5
	for i := 0; i < 10; i ++ {
		space.CreateShip(p, e2)
	}

	space.CreateShip(p, e3)

	Step(&space)
	if p.Control >= 0.5 {
		t.Errorf("Control did not decrease %s", p)
	}

	counter := 0
	for p.Empire == e1.Id {
		Step(&space)
		counter++
	}

	if p.Empire != e2.Id{
		t.Errorf("planet did not change owner %s", p)
	}

	if counter < 10 {
		t.Errorf("planet changed owner to fast %s %d", p ,counter)
	}
}
