package main

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/client/simple"
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	"fmt"
	"sync"
	"time"
)


var wg sync.WaitGroup

func main() {
	// find out how many empires are there
	c := client.Client{}
	space, err := c.PollState()

	if err != nil {
		panic(err)
	}

	for id := range space.Empires {
		wg.Add(1)
		go ControlLoop(id, simple.DistributeStrategy)
	}

	c.Close()

	wg.Wait()
}

func ControlLoop(empire uint32, f func(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command)) {
	defer func() {
		fmt.Println(empire, " lost")
		wg.Done()
	}()

	fmt.Println("Starting ", empire)
	c := client.Client{}
	defer c.Close()

	for {
		space, err := c.PollState()

		if err != nil {
			c.Close()
			continue
		}

		start := time.Now()
		response := pb.Command{
			Empire: empire,
		}

		owns_planet := false
		for _, planet := range space.Planets {
			if planet.Empire != empire {
				continue
			}

			f(space, planet, empire, &response)
			owns_planet = true
		}

		if ! owns_planet{
			break
		}

		fmt.Println(empire, time.Now().Sub(start))

		err = c.SendCommand(&response)
		if err != nil {
			fmt.Println(err)
		}
	}
}
