package strategies

import (
	commands2 "github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"math"

	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

type Distributed struct {
	matchingspace *state.Space
	graph         state.Graph
	empire        *state.Empire
}

func (d *Distributed) Init(empire *state.Empire) {
	d.empire = empire
}

func (dist *Distributed) Commands(space *state.Space) []commands2.Command {
	// cache graph
	if dist.matchingspace != space {
		// make sure we have a graph about this space
		dist.graph = state.NewGraph(space.Planets)
		dist.matchingspace = space
	}

	var commands []commands2.Command

	for _, myFleet := range dist.empire.Fleets {

		if myFleet.Position.Empire != dist.empire || myFleet.Position.Control < 0.5{
			continue
		}

		enemiesPresent := false
		for _, fleet := range myFleet.Position.Fleets{
			if fleet.Empire != dist.empire{
				enemiesPresent = true
				break
			}
		}

		if enemiesPresent{
			continue
		}

		// search weakest neighbor with enemy ships or planet
		lowestNeighborShips := math.MaxInt32
		var lowestNeighbor *state.Planet = nil

		for _, planet := range myFleet.Position.Connected {
			enemyShips := 0
			ownShips := 0
			for _, fleet := range planet.Fleets {
				if fleet.Empire != dist.empire {
					enemyShips += fleet.Size()
				} else {
					ownShips += fleet.Size()
				}
			}

			// enemies present
			if planet.Empire != dist.empire || enemyShips > 0 {

				// remember which neighbor is the weakest
				if ownShips == 0 && enemyShips < lowestNeighborShips {
					// that target has none of our ships.
					lowestNeighbor = planet
					lowestNeighborShips = enemyShips
				} else if lowestNeighborShips > enemyShips-ownShips {
					lowestNeighbor = planet
					lowestNeighborShips = enemyShips - ownShips
				}
			}
		}

		if lowestNeighborShips != math.MaxInt32 {
			// if we have more ships send them to the neighbor
			if myFleet.Size() > int(float32(lowestNeighborShips)*1.5) {
				// continue with next fleet
				break
			}

			commands = append(commands, commands2.MoveCommand{Destination: lowestNeighbor, Fleet: myFleet})
		} else {
			// there are no bordering enemy planets or fleets. Send myFleet in the direction of trouble

			var targetId *state.Planet
			// figure out the next planet with a enemy fleet or planet
			dist.graph.Visit(myFleet.Position, func(n state.Node) bool {
				if n.Planet.Empire != dist.empire {
					targetId = n.Planet
					return false
				}
				for _, fleet := range n.Planet.Fleets {
					if fleet.Empire != dist.empire {
						targetId = n.Planet
						return false
					}
				}
				return true
			})

			target := dist.graph.ShortestPath(myFleet.Position, targetId, true)
			if len(target) < 2 {
				return nil
			}

			commands = append(commands, commands2.MoveCommand{Destination: target[1], Fleet: myFleet})
		}
	}

	return commands
}
