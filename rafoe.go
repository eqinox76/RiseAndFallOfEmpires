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
)

type Client struct {
	output chan []byte
	Done   bool
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	space := state.NewSpace(5)
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
			fmt.Println("New client")
			mutex.Lock()

			client := Client{output: make(chan []byte)}
			workers = append(workers, &client)

			mutex.Unlock()
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
					fmt.Println("Client closed")
					removed++
				} else {
					w.output <- msg
				}
			}

			if removed > 0 {
				workers = workers[removed:]
			}

			mutex.Unlock()
		}
	}()

	// compute game state
	for {
		start := time.Now()
		engine.Step(&space)
		bytes, err := state.Serialize(&space)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		fanOut <- bytes
		fmt.Println(time.Now().Sub(start))

		// scan for one second the commands. After that compute the next state
		for start.Add(100 * time.Millisecond).After(time.Now()) {
			select {
			case cmd := <-commands:
				engine.ProcessCommand(&space, cmd)
			default:
				time.Sleep(10 * time.Millisecond)
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

		_, err := con.Write(msg)

		if err != nil {
			fmt.Errorf("%s throws %s", worker, err)
			break
		}
	}
	worker.Done = true
}
