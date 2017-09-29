package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	"math/rand"
	"math"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
)

func Step(space *state.Space) {
	space.Round++
	for _, planet := range space.Planets {

		computeControl(space, planet)
		computeFight(space, planet)
	}

	for _, empire := range space.Empires {
		// compute production for not passice empires
		// a empire can at most produce the
		if !empire.Passive {
			var totalControl float64 = 0.
			for planet, _ := range empire.Planets {
				totalControl += float64(space.Planets[planet].Control)
			}

			for planet, _ := range empire.Planets {
				if len(empire.Ships) > int(math.Sqrt(totalControl) * 100) {
					break
				}
				computeProduction(space, space.Planets[planet])
			}

		}
		if len(empire.Planets)+len(empire.Ships) == 0 {
			delete(space.Empires, empire.Id)
		}
	}
	return
}
func computeProduction(space *state.Space, planet *pb.Planet) {
	if rand.Float32() <= planet.Control {
		e := space.Empires[planet.Empire]

		space.CreateShip(planet, e)
	}
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

	// if no one or only the controlling empire is present control increases
	if len(planet.Orbiting) == 0 || (ownFleetPresent && len(fleets) == 1) {
		if planet.Control != 1 {
			if planet.Control > 0.999 {
				planet.Control = 1
			} else if planet.Control <= 0 {
				planet.Control = 0.001
			} else if planet.Control > 0.5 {
				planet.Control = planet.Control + (rand.Float32() * 0.1 * (1 - planet.Control))
			} else {
				planet.Control = planet.Control + (rand.Float32() * 0.1 * (planet.Control))
			}
		}
	} else {
		// else control decreases
		planet.Control -= float32(len(planet.Orbiting)) * 0.0005

		// if the control is to low and there is only one enemy fleet it takes over
		if planet.Control < 0 && len(fleets) == 1 {
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
