package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	space := state.NewSpace(9)
	eng := engine.GameEngine{Space: &space}
	eng.Init()

	// compute game state until the game is done
	for !space.Won() {
		start := time.Now()
		eng.Step()

		fmt.Println(time.Now().Sub(start), "for step,", len(space.Empires), "empires active,", len(space.Fleets), "Fleets")

		start = time.Now()
	}
}
