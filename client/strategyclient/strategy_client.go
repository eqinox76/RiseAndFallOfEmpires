package main

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/client/simple"
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	"fmt"
	"sync"
	_ "net/http/pprof"
	"log"
	"net/http"
	"github.com/eqinox76/RiseAndFallOfEmpires/client/special"
	"math/rand"
)

var wg sync.WaitGroup

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// find out how many empires are there
	c := client.Client{}
	space, err := c.PollState()

	clientChannels := make(map[uint32]chan *pb.Space, 0)
	// TODO need to close the channel when the empire has lost
	responseChannel := make(chan *pb.Command)
	doneChannel := make(chan uint32, 100)

	if err != nil {
		panic(err)
	}

	for id := range space.Empires {
		if id == 0 {
			continue
		}
		wg.Add(1)
		clientChannels[id] = make(chan *pb.Space, 20)
		ki := rand.Intn(3)
		switch ki {
		case 0:
			d := simple.Distributed{}
			go ControlLoop(id, d.DistributeStrategy, clientChannels[id], doneChannel, responseChannel)
		case 1:
			go ControlLoop(id, simple.RandomStrategy, clientChannels[id], doneChannel, responseChannel)
		case 2:
			go ControlLoop(id, special.FergsnStrategy, clientChannels[id], doneChannel, responseChannel)
		}
	}

	finishChannel := make(chan bool)

	defer func() {
		finishChannel <- true
		close(finishChannel)
		close(responseChannel)
	}()

	go func() {
		defer c.Close()
		// wait until the finish channel has been called to clean everything up
		for len(finishChannel) == 0 {
			space, err := c.PollState()

			if err != nil {
				c.Close()
				continue
			}

			for _, c := range clientChannels {
				c <- space
			}

			// check if a client is done
			select {
			case empire := <-doneChannel:
				fmt.Println(empire, " is done.")

				close(clientChannels[empire])
				delete(clientChannels, empire)
			default:
			}
		}
		for _, c := range clientChannels {
			close(c)
		}
	}()

	// send commands
	go func() {
		for response := range responseChannel {
			err := c.SendCommand(response)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	wg.Wait()
}

func ControlLoop(empire uint32, f func(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command), state_channel chan *pb.Space, doneChannel chan uint32, response_channel chan *pb.Command) {
	defer func() {
		doneChannel <- empire
		wg.Done()
	}()

	fmt.Println("Starting ", empire)

	for space := range state_channel {
		if len(state_channel) > 0 {
			// we have a never game state in the pipeline
			continue
		}

		// if we have lost terminate
		emp, ok := space.Empires[empire]
		if !ok || len(emp.Ships)+len(emp.Planets) == 0 {
			break
		}
		// or if the game has ended
		if len(space.Empires) <= 2{
			break
		}

		response := pb.Command{
			Empire: empire,
		}

		for _, planet := range space.Planets {
			if planet.Empire != empire {
				continue
			}

			f(space, planet, empire, &response)
		}

		response_channel <- &response
	}
}
