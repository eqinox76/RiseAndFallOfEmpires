package strategyclient

import (
	"fmt"
	"math/rand"

	"github.com/eqinox76/RiseAndFallOfEmpires/client/simple"
	"github.com/eqinox76/RiseAndFallOfEmpires/client/special"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
)

func RunClients(space *pb.Space, input chan *pb.Space, commands chan *pb.Command) {

	clientChannels := make(map[uint32]chan *pb.Space, 0)
	// TODO need to close the channel when the empire has lost
	responseChannel := make(chan *pb.Command)
	doneChannel := make(chan uint32, 100)

	for id := range space.Empires {
		if id == 0 {
			continue
		}
		clientChannels[id] = make(chan *pb.Space, 20)
		ki := rand.Intn(4)
		switch ki {
		case 0:
			fallthrough
		case 1:
			d := simple.Distributed{}
			go ControlLoop(id, d.DistributeStrategy, clientChannels[id], doneChannel, responseChannel)
		case 2:
			go ControlLoop(id, simple.RandomStrategy, clientChannels[id], doneChannel, responseChannel)
		case 3:
			go ControlLoop(id, special.FergsnStrategy, clientChannels[id], doneChannel, responseChannel)
		}
	}

	// send commands
	go func() {
		for response := range responseChannel {
			commands <- response
		}
	}()

	for space := range input {

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
}

func ControlLoop(empire uint32, f func(space *pb.Space, planet *pb.Planet, empire uint32, response *pb.Command), state_channel chan *pb.Space, doneChannel chan uint32, response_channel chan *pb.Command) {
	defer func() {
		doneChannel <- empire
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
		if len(space.Empires) <= 2 {
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
