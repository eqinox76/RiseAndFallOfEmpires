package engine

import (
	"math"
	"math/rand"

	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
)

func Step(space *state.Space) {
	space.Round++
	for _, planet := range space.Planets {

		computeControl(space, planet)
		computeFight(space, planet)
	}

	for _, empire := range space.Empires {
		// compute production for not passive empires

		if empire.Passive {
			continue
		}

		var totalControl float64 = 0.
		fullPlanets := make([]*pb.Planet, 0)

		for planet, _ := range empire.Planets {
			pl := space.Planets[planet]
			totalControl += float64(pl.Control)
			if pl.Control == 1 {
				fullPlanets = append(fullPlanets, pl)
			}
		}

		prod := math.Ceil(math.Log2(float64(len(fullPlanets))) + 0.1)

		// a empire can at most produce 100 divisions per fully controlled planet
		prod = math.Min(prod, (math.Sqrt(totalControl)*100)-float64(len(empire.Ships)))

		for ; prod > 0; prod-- {
			space.CreateShip(fullPlanets[rand.Intn(len(fullPlanets))], empire)
		}

		if len(empire.Planets)+len(empire.Ships) == 0 {
			delete(space.Empires, empire.Id)
		}
	}
	return
}

func ProcessCommand(space *state.Space, command *pb.Command) {
	for _, cmd := range command.Orders {
		switch order := cmd.Order.(type) {
		case *pb.Command_Order_Move:
			//validate input
			_, valid := space.Planets[order.Move.Start].Orbiting[order.Move.Ship]
			if valid {
				space.MoveShip(order.Move.Ship, order.Move.Start, order.Move.Destination)
			}
		}
	}
}

func computeFight(space *state.Space, planet *pb.Planet) {
	fleets := state.GetFleets(space.Ships, planet)

	if len(fleets) < 2 {
		// at most one empire present
		return
	}

	// fight all combinations
	for target, t_fleet := range fleets {
		lost := 0
		for attacker, a_fleet := range fleets {
			if target == attacker {
				continue
			}
			// let <ships>/<enemies> fight against that fleet
			lost += computeDamage(len(a_fleet)/len(fleets), len(t_fleet))
		}

		lost = util.MinInt(lost, len(t_fleet))
		for lost > 0 {
			space.RemoveShip(t_fleet[lost-1])
			lost--
		}
	}
}

func computeControl(space *state.Space, planet *pb.Planet) {
	fleets := state.GetFleets(space.Ships, planet)
	_, ownFleetPresent := fleets[planet.Empire]

	increase := func(control float32) float32 {
		if control > 0.999 {
			return 1
		} else if control <= 0 {
			return 0.001
		} else if control > 0.5 {
			return control + (rand.Float32() * 0.1 * (1 - control))
		} else {
			return control + (rand.Float32() * 0.1 * (control))
		}
	}
	// if no one or only the controlling empire is present control increases
	if len(planet.Orbiting) == 0 || (ownFleetPresent && len(fleets) == 1) {
		planet.Control = increase(planet.Control)
	} else {
		// adjust control slowly to reflect the fleet strengths
		targetControl := float32(len(fleets[planet.Empire])) / float32(len(planet.Orbiting))
		step := planet.Control - targetControl
		if step > 0 {
			step *= 0.1
			planet.Control -= step
		} else {
			planet.Control = increase(planet.Control)
		}

		// if the control is too low and there is only one enemy fleet it takes over
		if planet.Control < 0.05 && len(fleets) == 1 {
			planet.Control = 0
			for empire, _ := range fleets {
				delete(space.Empires[planet.Empire].Planets, planet.Id)
				planet.Empire = empire
				space.Empires[empire].Planets[planet.Id] = true
			}
		}
	}
}

func computeDamage(ships int, target int) int {
	const prob float64 = 0.05
	var deviation float64 = math.Sqrt(float64(ships) * prob * (1 - prob))
	var destroyed float64 = rand.NormFloat64()*deviation + (float64(ships) * prob)

	if destroyed < 0 {
		return 0
	} else if target < int(destroyed) {
		return target
	} else {
		return int(destroyed)
	}
}
