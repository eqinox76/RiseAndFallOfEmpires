package main

import (
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/eqinox76/RiseAndFallOfEmpires/client/strategyclient"
	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/mohae/deepcopy"
)

type gameServer struct {
	commands chan *pb.Command
	space    *state.Space
}

var maxWaitForCommands = flag.Int("maxWaitForCommands", 500, "max time we wait for all clients to send commands in ms")
var path = flag.String("write", "state.tmp", "were to write the data to")

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	stateChannel := make(chan *pb.Space)

	defer close(stateChannel)

	startPoint := state.NewSpace(9)
	server := gameServer{
		commands: make(chan *pb.Command),
		space:    &startPoint,
	}

	go strategyclient.RunClients(&server.space.Space, stateChannel, server.commands)

	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	writer, err := gzip.NewWriterLevel(f, gzip.DefaultCompression)
	if err != nil {
		panic(err)
	}

	defer writer.Close()

	// compute game state until the game is done
	for !server.space.Won() {
		start := time.Now()
		engine.Step(server.space)

		fmt.Println(time.Now().Sub(start), "for step", len(server.space.Empires), "empires active")
		start = time.Now()

		//we need to deepcopy state.Space because the deserialization by grpc and this step computation may interleave
		space, ok := deepcopy.Copy(server.space).(*state.Space)
		if !ok {
			panic(space)
		}
		//fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		stateChannel <- &server.space.Space

		data, err := space.Space.Marshal()
		err = binary.Write(writer, binary.LittleEndian, uint32(len(data)))
		if err != nil {
			panic(err)
		}
		_, err = writer.Write(data)
		if err != nil {
			panic(err)
		}

		fmt.Println(time.Now().Sub(start), "for writting to disk (", len(data), "bytes uncompressed)")
		start = time.Now()
		server.space = space

		commandsReceived := make(map[uint32]bool)
		// scan for the commands. After that compute the next state
		// check if all empires (minus the passive one) have send a command
		for len(commandsReceived) != len(server.space.Empires)-1 && start.Add(time.Duration(*maxWaitForCommands)*time.Millisecond).After(time.Now()) {
			select {
			case cmd := <-server.commands:
				engine.ProcessCommand(server.space, cmd)
				commandsReceived[cmd.Empire] = true
			default:
				time.Sleep(5 * time.Millisecond)
			}
		}

		fmt.Println(time.Now().Sub(start), "for processing commands")
		start = time.Now()
	}
}
