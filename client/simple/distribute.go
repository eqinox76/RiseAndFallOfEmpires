package simple

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math"
)

type Distributed struct {
	matchingspace *pb.Space
	graph         state.Graph
}

func (dist *Distributed) DistributeStrategy(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command) {
	if dist.matchingspace != space {
		// make sure we have a graph about this space
		dist.graph = state.NewGraph(space.Planets)
		dist.matchingspace = space
	}

	// search weakest neighbor with enemy ships or planet
	lowest_neighbor_ships := math.MaxInt32
	lowest_neighbor := uint32(0)
	fleets := state.GetFleets(space.Ships, planet)
	own_fleet, _ := fleets[empire]

	for _, p_id := range planet.Connected {
		fleets := state.GetFleets(space.Ships, space.Planets[p_id])
		own_fleet, ok := fleets[empire]

		if space.Planets[p_id].Empire != empire || len(fleets) > 1 || !ok {
			// enemies present

			// remember which neighbor is the weakest
			if !ok {
				// that target has none of our ships. send all there
				lowest_neighbor = p_id
				lowest_neighbor_ships = 0
				break
			} else if lowest_neighbor_ships > len(own_fleet) {
				lowest_neighbor_ships = len(own_fleet)
				lowest_neighbor = p_id
			}
		}
	}

	if lowest_neighbor_ships != math.MaxInt32 {
		// if we have more ships send them to the neighbor
		my_ships, ok := state.GetFleets(space.Ships, planet)[empire]
		if !ok || len(my_ships) < lowest_neighbor_ships {
			return
		}

		amount_to_sent := ((len(my_ships) + lowest_neighbor_ships) / 2) - lowest_neighbor_ships
		//fmt.Print(" my ships:", len(my_ships), " to send:", amount_to_sent)
		for count := 0; count < amount_to_sent; count ++ {
			order := pb.MovementOrder{
				Ship:        my_ships[count].Id,
				Start:       planet.Id,
				Destination: lowest_neighbor,
			}

			response.Orders = append(response.Orders, &pb.Command_Order{
				Order: &pb.Command_Order_Move{
					Move: &order,
				},
			})
		}
		return
	}

	if len(fleets) == 1 {
		// there are no bordering enemy planets. Send the whole fleet in the direction of trouble
		var target_id uint32
		dist.graph.Visit(planet.Id, func(n state.Node) bool {
			if n.Planet.Empire != empire {
				target_id = n.Planet.Id
				return false
			}
			return true
		})

		target := dist.graph.ShortestPath(planet.Id, target_id, true)
		if len(target) < 2{
			return
		}

		for count := 0; count < (len(own_fleet) * 2 / 3.); count ++ {
			order := pb.MovementOrder{
				Ship:        own_fleet[count].Id,
				Start:       planet.Id,
				Destination: target[1],
			}

			response.Orders = append(response.Orders, &pb.Command_Order{
				Order: &pb.Command_Order_Move{
					Move: &order,
				},
			})
		}
	}
}
