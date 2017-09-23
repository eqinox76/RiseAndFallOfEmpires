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
)

type clientConn struct {
	output chan *pb.Space
	Done   bool
}

type gameServer struct {
	clients []clientConn
	mutex   sync.Mutex
}

func (server *gameServer) Observe(id *pb.ID, client pb.GameServer_ObserveServer) error {
	log.Println("Client",client.Context(),"connected")
	// register us for the game states
	channel := make(chan *pb.Space)
	c := clientConn{channel, false}
	server.mutex.Lock()
	server.clients = append(server.clients, c)
	server.mutex.Unlock()

	defer func() { c.Done = true }()
	for space := range channel {
		err := client.Send(space)
		if err != nil {
			return err
		}
	}
	log.Println(len(channel), "The channel for this client has been closed")
	return nil
}

func (server *gameServer) StrategyClient(stream pb.GameServer_StrategyClientServer) error {
	return nil
}

func (server *gameServer) fanOut(input chan *pb.Space) {
	for msg := range input {

		removed := 0
		server.mutex.Lock()
		for i, client := range server.clients {
			if client.Done {
				// this worker is done so close it
				close(client.output)
				server.clients[i] = server.clients[removed]
				removed++
			} else {
				client.output <- msg
				//fmt.Println(w, len(w.output), len(fanOut))
			}
		}

		if removed > 0 {
			server.clients = server.clients[removed:]
		}

		server.mutex.Unlock()
	}
}

var maxWaitForCommands = flag.Int("maxWaitForCommands", 1000, "max time we wait for all clients to send commands in ms")

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	space := state.NewSpace(9)
	fanOut := make(chan *pb.Space)
	commands := make(chan *pb.Command)
	defer close(fanOut)

	lis, err := net.Listen("tcp", ":9076")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// TODO do we need to secure the local connection?
	grpc.WithInsecure()

	grpcServer := grpc.NewServer()
	server := gameServer{}
	pb.RegisterGameServerServer(grpcServer, &server)
	go func() {
		grpcServer.Serve(lis)
	}()

	// start fan out
	go server.fanOut(fanOut)

	// compute game state until the game is done
	for ! space.Won() {
		start := time.Now()
		engine.Step(&space)

		//fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		fanOut <- &space.Space
		fmt.Println(time.Now().Sub(start))

		commandsReceived := make(map[uint32]bool)
		// scan for the commands. After that compute the next state
		for len(commandsReceived) != len(space.Empires)-1 && start.Add(time.Duration(*maxWaitForCommands) * time.Millisecond).After(time.Now()) {
			select {
			case cmd := <-commands:
				engine.ProcessCommand(&space, cmd)
				commandsReceived[cmd.Empire] = true
				// check if all empires (minus the passive one) have send a command
			default:
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
}

//func process_commands(conn net.Conn, client *Client, commands chan *pb.Command) {
//
//	for {
//		header := make([]byte, 4)
//		_, err := conn.Read(header)
//
//		if err != nil {
//			fmt.Println(err)
//			conn.Close()
//			client.Done = true
//			return
//		}
//
//		l := binary.LittleEndian.Uint32(header)
//
//		msgbuffer := make([]byte, l)
//		_, err = conn.Read(msgbuffer)
//		if err != nil {
//			fmt.Println(err)
//			conn.Close()
//			client.Done = true
//			return
//		}
//
//		command := pb.Command{}
//		err = proto.Unmarshal(msgbuffer, &command)
//		if err != nil {
//			fmt.Println(err)
//			conn.Close()
//			client.Done = true
//			return
//		}
//
//		commands <- &command
//	}
//}
//
//func forward_game_state(con net.Conn, worker *Client) {
//	for msg := range worker.output {
//
//		con.SetWriteDeadline(time.Now().Add(150 * time.Millisecond))
//		_, err := con.Write(msg)
//
//		if err != nil {
//			fmt.Println(worker, " closed due to ", err)
//			break
//		}
//	}
//	worker.Done = true
//}
