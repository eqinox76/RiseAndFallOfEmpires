package main

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"math/rand"
	"fmt"
	"sync"
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
	defer func(){
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

			if rand.Float32() > 5. / float32(len(space.Empires[empire].Planets)){
				//only send commands for some planets
				continue
			}

			fleets := state.GetFleets(space.Ships, planet)
			my_fleet := fleets[empire]
			// if there are no enemies and we have more than 10 ships
			if len(fleets) == 1 && len(my_fleet) > 10 {
				// in some cases
				if rand.Float32() < float32(len(my_fleet)) / 100. {
					// send random amount of ships to random target
					sent := rand.Intn(len(my_fleet))
					// make sure at least one ship is sent
					sent |= 1
					for i := 0; i < sent; i++ {
						target := planet.Connected[rand.Intn(len(planet.Connected))]
						order := pb.MovementOrder{
							Ship:        my_fleet[i].Id,
							Start:       planet.Id,
							Destination: target,
						}

						response.Orders = append(response.Orders, &pb.Command_Order{
							Order: &pb.Command_Order_Move{
								Move: &order,
							},
						})
					}

					err := c.SendCommand(&response)
					if err != nil{
						fmt.Println(err)
					}					else {
					fmt.Println("Send commands for empire", empire)
					}
				}
			}
		}
	}
}
