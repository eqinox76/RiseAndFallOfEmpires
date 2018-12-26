package strategies

import (
	commands2 "github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
)

type RandomStrategy struct {
	empire *state.Empire
}

func (strategy *RandomStrategy) Init(empire *state.Empire) {
	strategy.empire = empire
}

func (strategy RandomStrategy) Commands(space *state.Space) []commands2.Command {
	if rand.Float32() > 5./float32(len(strategy.empire.Planets)) {
		//only send commands for some planets
		return nil
	}

	var commands []commands2.Command

	for _, fleet := range strategy.empire.Fleets {
		onlyOwn := true
		// if there are enemies we do not move the fleet
		for _, orbiting := range fleet.Position.Fleets {
			if orbiting.Empire != strategy.empire {
				onlyOwn = false
				break
			}
		}

		if !onlyOwn {
			continue
		}

		// need at least 10 ships in the fleet
		if fleet.Size() > 10 {
			// in some cases
			if rand.Float32() < float32(fleet.Size())/100. {
				// send fleet to random target
				dest := fleet.Position.Connected[rand.Intn(len(fleet.Position.Connected))]
				commands = append(commands, commands2.MoveCommand{Destination: dest, Fleet: fleet})
			}
		}
	}

	return commands
}
