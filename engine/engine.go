package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/strategies"
	"github.com/eqinox76/RiseAndFallOfEmpires/util"
	"log"
	"math"
	"math/rand"
)

type GameEngine struct {
	Space  *state.Space
	strats []strategies.Strategy
}

func (engine *GameEngine) Init() {
	engine.strats = make([]strategies.Strategy, len(engine.Space.Empires))

	for id, empire := range engine.Space.Empires {
		if empire.Passive {
			engine.strats[id] = &strategies.NoneStrategy{}
			continue
		}

		ki := rand.Intn(2)
		switch ki {
		case 0:
			engine.strats[id] = &strategies.Distributed{}
		case 1:
			engine.strats[id] = &strategies.RandomStrategy{}
		}
		engine.strats[id].Init(empire)
	}
}

func (engine *GameEngine) Step() {
	engine.Space.Round++
	for _, planet := range engine.Space.Planets {

		computeControl(planet)
		engine.computeFight(planet)
	}

	for i := len(engine.Space.Empires) - 1; i >= 0; i-- {
		if len(engine.Space.Empires[i].Planets)+len(engine.Space.Empires[i].Fleets) == 0 {
			// del empire
			engine.Space.Empires[i] = engine.Space.Empires[len(engine.Space.Empires)-1]
			engine.Space.Empires = engine.Space.Empires[:len(engine.Space.Empires)-1]
		}
	}

	// production
	for _, empire := range engine.Space.Empires {
		if empire.Passive {
			continue
		}

		totalControl := 0.
		fullPlanets := make([]*state.Planet, 0)

		for _, planet := range empire.Planets {
			totalControl += float64(planet.Control)
			if planet.Control == 1 {
				fullPlanets = append(fullPlanets, planet)
			}
		}

		prod := math.Ceil(math.Log2(float64(len(fullPlanets))) + 0.1)

		fleetSize := 0.
		for _, fleet := range empire.Fleets {
			fleetSize += float64(fleet.Size())
		}

		// a empire can at most produce sqrt(100 divisions per fully controlled planet)
		prod = math.Min(prod, (math.Sqrt(totalControl)*100)-fleetSize)

		for ; prod > 0; prod-- {
			pl := fullPlanets[rand.Intn(len(fullPlanets))]
			if pl.Production > 0.9 {

				fleet := pl.EmpireFleet(empire)
				if fleet == nil {
					fleet = engine.Space.CreateFleet(pl, empire)
				}

				switch rand.Int31n(3) {
				case 0:
					fleet.LightSquads++
				case 1:
					fleet.HeavySquads++
				case 2:
					fleet.RangedSquads++
				}

				pl.Production = 0
			} else {
				pl.Production += 0.1
			}
		}

	}

	cmds := make([]commands.Command, 0)
	for _, s := range engine.strats {
		cmds = append(cmds, s.Commands(engine.Space)...)
	}

	for _, cmd := range cmds {
		if ! cmd.Validate() {
			log.Fatal(cmd, "Invalid command")
		} else {
			cmd.Execute(engine.Space)
		}
	}
	return
}

func (engine *GameEngine) computeFight(planet *state.Planet) {
	fleets := make(map[*state.Empire]*state.Fleet)
	allSize := 0

	for _, fleet := range planet.Fleets {
		fleets[fleet.Empire] = fleet
		allSize += fleet.Size()
	}

	if len(fleets) <= 1 {
		// at most one empire present
		return
	}

	// fight all combinations
	for target, tFleet := range fleets {
		lost := 0
		for attacker, aFleet := range fleets {
			if target == attacker {
				continue
			}
			// let <ships>/<enemies> fight against that fleet
			lost += computeDamage(aFleet.Size()/(len(fleets)-1), tFleet.Size())
		}

		lost = util.MinInt(lost, tFleet.Size())
		for lost > 0 {
			// TODO find a better way to destroy a random ship
			if tFleet.LightSquads > 0 && tFleet.HeavySquads > 0 && tFleet.RangedSquads > 0 {
				switch rand.Intn(3) {
				case 0:
					tFleet.LightSquads--
				case 1:
					tFleet.HeavySquads--
				case 2:
					tFleet.RangedSquads--
				}
			} else if tFleet.LightSquads > 0 && tFleet.HeavySquads > 0 {

				switch rand.Intn(2) {
				case 0:
					tFleet.LightSquads--
				case 1:
					tFleet.HeavySquads--
				}
			} else if tFleet.LightSquads > 0 && tFleet.RangedSquads > 0 {

				switch rand.Intn(2) {
				case 0:
					tFleet.LightSquads--
				case 1:
					tFleet.RangedSquads--
				}
			} else if tFleet.HeavySquads > 0 && tFleet.RangedSquads > 0 {

				switch rand.Intn(2) {
				case 0:
					tFleet.HeavySquads--
				case 1:
					tFleet.RangedSquads--
				}
			} else if tFleet.LightSquads > 0 {
				tFleet.LightSquads -= lost
				lost = 0
			} else if tFleet.HeavySquads > 0 {
				tFleet.HeavySquads -= lost
				lost = 0
			} else if tFleet.RangedSquads > 0 {
				tFleet.RangedSquads -= lost
				lost = 0
			}

			lost--
		}
	}

	// check which fleets did not make it
	for i := len(planet.Fleets) - 1; i >= 0; i-- {
		if planet.Fleets[i].Size() == 0 {
			engine.Space.DestroyFleet(planet.Fleets[i])
		}
	}
}

func computeControl(planet *state.Planet) {
	fleets := make(map[*state.Empire]int)

	for _, fleet := range planet.Fleets {
		fleets[fleet.Empire] += fleet.Size()
	}

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
	if len(planet.Fleets) == 0 || (ownFleetPresent && len(fleets) == 1) {
		planet.Control = increase(planet.Control)
	} else {
		// adjust control slowly to reflect the fleet strengths
		targetControl := float32(fleets[planet.Empire]) / float32(len(planet.Fleets))
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
			planet.Production = 0
			for empire := range fleets {
				for i := range planet.Empire.Planets {
					l := len(planet.Empire.Planets) - 1
					planet.Empire.Planets[i] = planet.Empire.Planets[l]
					planet.Empire.Planets = planet.Empire.Planets[:l]
					break
				}

				planet.Empire = empire
				empire.Planets = append(empire.Planets, planet)
				break
			}
		}
	}
}

func computeDamage(ships int, target int) int {
	const prob float64 = 0.2
	var deviation = math.Sqrt(float64(ships) * prob * (1 - prob))
	var destroyed = rand.NormFloat64()*deviation + (float64(ships) * prob)

	if destroyed < 0 {
		return 0
	} else if target < int(destroyed) {
		return target
	} else {
		return int(destroyed)
	}
}
