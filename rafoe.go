package main

import (
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/model"
	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
)

func main() {

	space := model.NewSpace()

	for i:=0 ; i< 8 ; i++ {
		bytes, _ := space.Serialize()
		fmt.Printf("serialize: %d, Planets: %d Ships: %d\n", len(bytes), len(space.Planets), len(space.Ships))
		engine.Step(&space)
	}

}
