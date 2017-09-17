package special

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
)

// this function will be called for all owned planets
func FergsnStrategy(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command) {
	// the orbiting fleets by empire
	fleets := state.GetFleets(space.Ships, planet)
	// get our own fleet
	my_fleet := fleets[empire]
	// if there are more than 10 of our own ships
	if len(my_fleet) > 10 {
		// target a random connected planet
		target := planet.Connected[rand.Intn(len(planet.Connected))]
		// and send 5 ships
		move(space, empire, response, 5, planet, target)
	}
}

func move(space *pb.Space, empire uint32, response *pb.Command, amount_ships uint64, from *pb.Planet, to uint32) {
	fleets := state.GetFleets(space.Ships, from)
	my_fleet := fleets[empire]
	for i := uint64(0); i < amount_ships; i++ {
		order := pb.MovementOrder{
			Ship:        my_fleet[i].Id,
			Start:       from.Id,
			Destination: to,
		}

		response.Orders = append(response.Orders, &pb.Command_Order{
			Order: &pb.Command_Order_Move{
				Move: &order,
			},
		})
	}
}
