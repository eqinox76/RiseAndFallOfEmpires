package main

import (
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
	"net"
	"sync"
	"time"
	"math/rand"
	"google.golang.org/grpc"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"flag"
	"log"
	"context"
	"github.com/mohae/deepcopy"
)

type clientConn struct {
	output chan *pb.Space
	Done   bool
}

type gameServer struct {
	clients  []clientConn
	mutex    sync.Mutex
	commands chan *pb.Command
	space    *state.Space
}

func (server *gameServer) addClient() *clientConn {
	channel := make(chan *pb.Space)
	c := clientConn{channel, false}
	server.mutex.Lock()
	server.clients = append(server.clients, c)
	server.mutex.Unlock()
	return &c
}

func (server *gameServer) Observe(id *pb.ID, client pb.GameServer_ObserveServer) error {
	log.Println("Client", client.Context(), "connected")
	// register us for the game states
	c := server.addClient()

	defer func() { c.Done = true }()

	for space := range c.output {
		err := client.Send(space)
		if err != nil {
			return err
		}
	}
	log.Println("The channel for this client has been closed")
	return nil
}

func (server *gameServer) CurrentGameState(ctx context.Context, id *pb.ID) (*pb.Space, error) {
	return &server.space.Space, nil
}

func (server *gameServer) StrategyClient(stream pb.GameServer_StrategyClientServer) error {
	log.Println("Strategy", stream.Context(), "connected")
	// register us for the game states
	c := server.addClient()

	go func() {
		for {
			cmd, err := stream.Recv()
			if err != nil {
				c.Done = true
				break
			}
			server.commands <- cmd
		}
	}()

	for space := range c.output {
		err := stream.Send(space)
		if err != nil {
			return err
		}
	}

	return nil
}

func (server *gameServer) fanOut(input chan *pb.Space) {
	for msg := range input {

		removed := 0
		server.mutex.Lock()
		for i, c := range server.clients {
			if c.Done {
				// this worker is done so close it
				close(c.output)
				server.clients[i] = server.clients[removed]
				removed++
			} else {
				c.output <- msg
				//fmt.Println(w, len(w.output), len(fanOut))
			}
		}

		if removed > 0 {
			server.clients = server.clients[removed:]
		}

		server.mutex.Unlock()
	}
}

var maxWaitForCommands = flag.Int("maxWaitForCommands", 2000, "max time we wait for all clients to send commands in ms")

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	fanOut := make(chan *pb.Space)

	defer close(fanOut)

	lis, err := net.Listen("tcp", ":9076")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	startPoint := state.NewSpace(9)
	server := gameServer{
		commands: make(chan *pb.Command),
		space:    &startPoint,
	}
	pb.RegisterGameServerServer(grpcServer, &server)
	go func() {
		grpcServer.Serve(lis)
	}()

	// start fan out
	go server.fanOut(fanOut)

	// compute game state until the game is done
	for ! server.space.Won() {
		start := time.Now()
		engine.Step(server.space)

		//fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		fanOut <- &server.space.Space
		fmt.Println(time.Now().Sub(start))

		//we need to deepcopy state.Space because the deserialization by grpc and this step computation may interleave
		space, ok := deepcopy.Copy(server.space).(*state.Space)
		if !ok{
			panic(space)
		}
		server.space = space

		commandsReceived := make(map[uint32]bool)
		// scan for the commands. After that compute the next state
		// check if all empires (minus the passive one) have send a command
		for len(commandsReceived) != len(server.space.Empires)-1 && start.Add(time.Duration(*maxWaitForCommands) * time.Millisecond).After(time.Now()) {
			select {
			case cmd := <-server.commands:
				engine.ProcessCommand(server.space, cmd)
				commandsReceived[cmd.Empire] = true
			default:
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
}
