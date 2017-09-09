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
)

type Worker struct {
	input chan []byte
	Done  bool
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	space := state.NewSpace(13)
	fanOut := make(chan []byte)
	defer close(fanOut)

	go func() {
		workers := make([]*Worker, 0)
		mutex := sync.Mutex{}

		go func() {
			// Listen for incoming connections.
			l, err := net.Listen("tcp", "localhost:9076")
			if err != nil {
				fmt.Println("Error listening:", err.Error())
				os.Exit(1)
			}

			defer l.Close()

			fmt.Println("Listening on localhost:9076")
			for {
				// Listen for an incoming connection.
				con, err := l.Accept()
				if err != nil {
					fmt.Println("Error accepting: ", err.Error())
					os.Exit(1)
				}

				mutex.Lock()

				w := Worker{input: make(chan []byte)}
				workers = append(workers, &w)

				mutex.Unlock()
				go send(con, &w)
			}
		}()

		for msg := range fanOut {
			mutex.Lock()
			for i, w := range workers {
				fmt.Println(w)
				if w.Done {
					// this worker is done so remove it

					close(w.input)
					workers[i] = workers[ len(workers)-1]
					workers = workers[:len(workers)-1]
				} else {
					w.input <- msg
				}
			}
			mutex.Unlock()
		}

	}()

	for {
		bytes, _ := state.Serialize(&space)
		fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		engine.Step(&space)
		time.Sleep(1000 * time.Millisecond)
		fanOut <- bytes
	}

}

func send(con net.Conn, worker *Worker) {
	for msg := range worker.input {

		_, err := con.Write(msg)

		if err != nil{
			break
		}
	}
	worker.Done = true
}
