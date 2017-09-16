package main

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/client/simple"
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	"fmt"
	"sync"
	"time"
	_ "net/http/pprof"
	"log"
	"net/http"
)

var wg sync.WaitGroup

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// find out how many empires are there
	c := client.Client{}
	space, err := c.PollState()

	state_channel := make(chan *pb.Space)
	response_channel := make (chan *pb.Command)

	if err != nil {
		panic(err)
	}

	for id := range space.Empires {
		if id == 0 {
			continue
		}
		wg.Add(1)
		if id % 2 == 0 {
			go ControlLoop(id, simple.DistributeStrategy, state_channel, response_channel)
		} else {
			go ControlLoop(id, simple.RandomStrategy, state_channel, response_channel)
		}
	}

	finish_channel := make(chan bool)

	defer func() {
		finish_channel <- true
		close(finish_channel)
		close(response_channel)
	}()

	go func() {
		defer close(state_channel)
		defer c.Close()
		for len(finish_channel) == 0 {
			space, err := c.PollState()

			if err != nil {
				c.Close()
				continue
			}
			state_channel <- space
		}
	}()
	go func() {
		for response := range response_channel {
			err := c.SendCommand(response)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()



	wg.Wait()
}

func ControlLoop(empire uint32, f func(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command), state_channel chan *pb.Space, response_channel chan *pb.Command) {
	defer func() {
		fmt.Println(empire, " lost")
		wg.Done()
	}()

	fmt.Println("Starting ", empire)

	for space := range state_channel {
		if len(state_channel) > 0 {
			// we have a never game state in the pipeline
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

		if ! owns_planet {
			break
		}

		fmt.Println(empire, time.Now().Sub(start))

		response_channel <- &response
	}
}
