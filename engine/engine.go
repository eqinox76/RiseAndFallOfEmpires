package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	"math/rand"
	"math"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
)

func Step(space *state.Space) {
	for _, planet := range space.Planets {
		if space.Empires[planet.Empire].Passive {
			// this empire produces nothing
			continue
		}

		if rand.Float32() <= planet.Control {
			// check if the empire can have a new ship
			e := space.Empires[planet.Empire]

			if (len(e.Planets) * 100) > len(e.Ships) {
				space.CreateShip(planet, e)
			}
		}

		computeControl(planet)

		computeFight(space, planet)

		computeOwner(space, planet)
	}
	return
}
func computeOwner(space *state.Space, planet *pb.Planet) {
	fleets := getFleets(space.Ships, planet)
	_ , fleet_exists := fleets[planet.Empire]
	if fleet_exists{
		// this planet has a defending fleet
		return
	}

	if len(fleets) > 1{
		// this planet is still being fought over
		return
	}

	for id := range fleets {
		old_owner := space.Empires[planet.Id]
		util.RemoveUint32(&old_owner.Planets, planet.Id)
		new_owner := space.Empires[id]
		new_owner.Planets = append(new_owner.Planets, planet.Id)
		planet.Empire = id
		planet.Control = 0.
	}
}

func getFleets(global_ships []*pb.Ship, planet *pb.Planet) map[uint32][]*pb.Ship {
	fleets := make(map[uint32][]*pb.Ship)

	for _, id := range planet.Orbiting {
		s := global_ships[id]
		_, ok := fleets[s.Empire]
		if !ok {
			fleets[s.Empire] = []*pb.Ship{}
		}

		fleets[s.Empire] = append(fleets[s.Empire], s)
	}

	return fleets
}

func computeFight(space *state.Space, planet *pb.Planet) {
	fleets := getFleets(space.Ships, planet)

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

func computeControl(planet *pb.Planet) {
	if planet.Control != 1 {
		if planet.Control > 0.999 {
			planet.Control = 1
		} else if planet.Control > 0.5 {
			planet.Control = planet.Control + (rand.Float32() * 0.1 * (1 - planet.Control))
		} else {
			planet.Control = planet.Control + (rand.Float32() * 0.1 * (planet.Control))
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
