package main

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"fmt"
	"sync"
	"math"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	go ControlLoop(1)

	wg.Add(1)
	go ControlLoop(2)
	wg.Add(1)
	go ControlLoop(3)
	wg.Add(1)
	go ControlLoop(4)
	wg.Add(1)
	go ControlLoop(5)

	wg.Wait()
}

func ControlLoop(empire uint32) {
	defer func() {
		fmt.Println("Done")
		wg.Done()
	}()

	c := client.Client{}
	for {
		space, err := c.PollState()

		if err != nil {
			panic(err)
		}

		response := pb.Command{
			Empire: empire,
		}

		for _, planet := range space.Planets {
			if planet.Empire != empire {
				continue
			}

			lowest_neighbor_ships := math.MaxInt32
			lowest_neighbor := uint32(0)
			for _, p_id := range planet.Connected {
				list, ok := state.GetFleets(space.Ships, space.Planets[p_id])[empire]
				if ok {
					// remember which neighbor is the weakest
					if lowest_neighbor_ships > len(list) {
						lowest_neighbor_ships = len(list)
						lowest_neighbor = p_id
					}
				}else{
					lowest_neighbor = p_id
					lowest_neighbor_ships = 0
					break
				}
			}

			// if we have more send them to the neighbor
			my_ships, ok := state.GetFleets(space.Ships, planet)[empire]
			if !ok || len(my_ships) < lowest_neighbor_ships {
				continue
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
			//fmt.Println()
		}

		err = c.SendCommand(&response)
		if err != nil {
			fmt.Println(err)
		}
	}
}
