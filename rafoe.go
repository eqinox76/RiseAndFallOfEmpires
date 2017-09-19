package main

import (
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
	"net"
	"os"
	"sync"
	"time"
	"math/rand"
	"github.com/golang/protobuf/proto"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"encoding/binary"
	"flag"
)

type Client struct {
	output chan []byte
	Done   bool
}

var maxWaitForCommands = flag.Int("maxWaitForCommands", 1000, "max time we wait for all clients to send commands in ms")

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	space := state.NewSpace(9)
	fanOut := make(chan []byte)
	commands := make(chan *pb.Command)
	defer close(fanOut)

	workers := make([]*Client, 0)
	mutex := sync.Mutex{}

	// start server port listener
	go func() {
		// Listen for incoming connections.
		l, err := net.Listen("tcp", "localhost:9076")
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}

		defer l.Close()
		fmt.Println("Listening on localhost:9076")

		// forever accept new connections
		for {
			con, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			mutex.Lock()

			client := Client{output: make(chan []byte)}
			workers = append(workers, &client)

			mutex.Unlock()
			fmt.Println("New client", client)
			// and start the client
			go forward_game_state(con, &client)
			go process_commands(con, &client, commands)
		}
	}()

	// start fan out
	go func() {
		for msg := range fanOut {

			removed := 0
			mutex.Lock()
			for i, w := range workers {

				if w.Done {
					// this worker is done so close it
					close(w.output)
					workers[i] = workers[removed]
					fmt.Println(w, " removed")
					removed++
				} else {
					w.output <- msg
					//fmt.Println(w, len(w.output), len(fanOut))
				}
			}

			if removed > 0 {
				workers = workers[removed:]
			}

			mutex.Unlock()
		}
	}()

	// compute game state forever
	for ! space.Won() {
		start := time.Now()
		engine.Step(&space)

		bytes, err := state.Serialize(&space)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		fanOut <- bytes
		fmt.Println(time.Now().Sub(start))

		commandsReceived := make(map[uint32]bool)
		// scan for the commands. After that compute the next state
		for len(commandsReceived) != len(space.Empires) - 1 && start.Add(time.Duration(*maxWaitForCommands) * time.Millisecond).After(time.Now()) {
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
func process_commands(conn net.Conn, client *Client, commands chan *pb.Command) {

	for {
		header := make([]byte, 4)
		_, err := conn.Read(header)

		if err != nil {
			fmt.Println(err)
			conn.Close()
			client.Done = true
			return
		}

		l := binary.LittleEndian.Uint32(header)

		msgbuffer := make([]byte, l)
		_, err = conn.Read(msgbuffer)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			client.Done = true
			return
		}

		command := pb.Command{}
		err = proto.Unmarshal(msgbuffer, &command)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			client.Done = true
			return
		}

		commands <- &command
	}
}

func forward_game_state(con net.Conn, worker *Client) {
	for msg := range worker.output {

		con.SetWriteDeadline(time.Now().Add(150 * time.Millisecond))
		_, err := con.Write(msg)

		if err != nil {
			fmt.Println(worker, " closed due to ", err)
			break
		}
	}
	worker.Done = true
}
